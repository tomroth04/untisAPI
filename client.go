package untisApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-resty/resty/v2"
	"github.com/pquerna/otp/totp"
	"github.com/rotisserie/eris"
	"github.com/tidwall/gjson"
	. "github.com/tomroth04/untisAPI/types"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO: add more information to the errors

type Client struct {
	BaseUrl            string
	School             string
	Identity           string
	Username           string
	Secret             string
	token              string
	sessionInformation SessionInformation
	httpClient         *resty.Client
}

func NewClient(baseUrl string, school string, identity string, username string, secret string) Client {
	return Client{
		BaseUrl:    "https://" + baseUrl,
		School:     school,
		Identity:   identity,
		Username:   username,
		Secret:     secret,
		httpClient: resty.New(),
	}
}

// Login with your credentials
//
// Notice: The server may revoke this session after less than 10min of idle.**
func (c *Client) Login() error {
	if err := c.getAccessToken(); err != nil {
		log.Println(err)
		return err
	}

	// Get personId & personType
	resp, err := c.httpClient.R().SetHeaders(
		c.getHeaders(),
	).Get(
		fmt.Sprintf("%s/WebUntis/api/app/config", c.BaseUrl),
	)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(fmt.Sprintf("status code non 200, body: %s", resp.String()))
	}

	if !gjson.GetBytes(resp.Body(), "data").Exists() {
		return eris.New("response didn't contain any data key")
	}

	if !gjson.GetBytes(resp.Body(), "data.loginServiceConfig").Exists() {
		return eris.New("response didn't contain any loginServiceConfig key")
	}

	if !gjson.GetBytes(resp.Body(), "data.loginServiceConfig.user").Exists() {
		return eris.New("response didn't contain any loginServiceConfig.user key")
	}

	if !gjson.GetBytes(resp.Body(), "data.loginServiceConfig.user.personId").Exists() {
		return eris.New("response didn't contain any personId")
	}

	c.sessionInformation.PersonId = int(gjson.Get(resp.String(), "data.loginServiceConfig.user.personId").Int())

	persons := gjson.Get(resp.String(), "data.loginServiceConfig.user.persons")
	if !persons.IsArray() {
		return eris.New("response didn't contain any persons array")
	}

	person := persons.Array()[0]
	c.sessionInformation.PersonType = int(person.Get("type").Int())

	resp, err = c.httpClient.R().SetHeader(
		"Cookie", c.GetCookie(),
	).Get(fmt.Sprintf("%s/WebUntis/api/daytimetable/config", c.BaseUrl))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 {
		if res := gjson.Get(resp.String(), "data.klasseId"); res.Type != gjson.Number {
			return nil
		} else {
			c.sessionInformation.ClassId = int(res.Int())
		}

	} else {
		return errors.New(fmt.Sprintf("Status code non 200, %s", resp.String()))
	}

	return nil
}

func (c *Client) Logout() error {
	resp, err := c.httpClient.R().SetQueryParam(
		"school", c.School,
	).SetBody(
		map[string]any{
			"id":      c.Identity,
			"method":  "logout",
			"params":  "{}",
			"jsonrpc": "2.0",
		}).SetHeaders(
		c.getHeaders(),
	).Post(
		fmt.Sprintf("%s/WebUntis/jsonrpc.do", c.BaseUrl))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New("status code non 200")
	}
	return nil
}

// Make JSON-RPC requests with the current session
func (c *Client) request(method string, params interface{}, validateSession bool) ([]byte, error) {
	if validateSession {
		if err := c.validateSession(); err != nil {
			return nil, err
		}
	}

	var resp *resty.Response

	if err := backoff.Retry(func() error {
		var err error
		resp, err = c.httpClient.R().SetQueryParam(
			"school", c.School,
		).SetHeader(
			"Cookie", c.GetCookie(),
		).SetBody(
			map[string]any{
				"id":      c.Identity,
				"method":  method,
				"params":  params,
				"jsonrpc": "2.0",
			},
		).Post(c.BaseUrl + "/WebUntis/jsonrpc.do")

		if err != nil {
			return err
		}
		if resp.IsError() {
			return errors.New("server response non 200e")
		}
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}

	if !gjson.Get(resp.String(), "result").Exists() {
		return nil, errors.New("server didn't return any result")
	}

	return resp.Body(), nil
}

func (c *Client) validateSession() error {
	b := backoff.NewExponentialBackOff()
	b.Multiplier = 3
	b.MaxElapsedTime = 30 * time.Minute
	return backoff.Retry(func() error {
		_, err := c.GetLatestSchoolyear(false)
		if err != nil {
			f := backoff.NewExponentialBackOff()
			b.MaxInterval = 30 * time.Minute
			b.Multiplier = 2.5
			return backoff.Retry(func() error {
				return c.Login()
			}, f)
		}
		return err
	}, b)
}

func (c *Client) requestTimeTable(id int, timeTableType int, startDate time.Time, endDate time.Time, validateSession bool) ([]GenericLesson, error) {
	params := map[string]any{
		"options": map[string]any{
			"id": time.Now().UnixMilli(),
			"element": map[string]any{
				"id":   id,
				"type": timeTableType,
			},
			"showLsText":       true,
			"showStudentgroup": true,
			"showLsNumber":     true,
			"showSubstText":    true,
			"showInfo":         true,
			"showBooking":      true,
			"klasseFields": []string{
				"id", "name", "longname", "externalkey",
			},
			"roomFields": []string{
				"id", "name", "longname", "externalkey",
			},
			"subjectFields": []string{
				"id", "name", "longname", "externalkey",
			},
			"teacherFields": []string{
				"id", "name", "longname", "externalkey",
			},
		},
	}

	if !startDate.IsZero() {
		params["options"].(map[string]any)["startDate"] = GetDateUntisFormat(startDate)
	}
	if !endDate.IsZero() {
		params["options"].(map[string]any)["endDate"] = GetDateUntisFormat(endDate)
	}

	resp, err := c.request("getTimetable", params, validateSession)
	if err != nil {
		return nil, err
	}

	res := gjson.GetBytes(resp, "result")
	if !res.Exists() {
		return nil, errors.New("no result in response")
	}

	return TransformResultLesson(res.Array()), nil
}

func (c *Client) GetTimetableForToday(id int, timeTableType int, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(id, timeTableType, time.Time{}, time.Time{}, validateSession)
}

func (c *Client) GetOwnTimetableForToday(validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(
		c.sessionInformation.PersonId,
		c.sessionInformation.PersonType,
		time.Time{}, time.Time{}, validateSession)
}

func (c *Client) GetTimetableFor(id int, timeTableType int, date time.Time, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(id, timeTableType, date, date, validateSession)
}

func (c *Client) GetOwnTimetableForRange(startDate time.Time, endDate time.Time, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(c.sessionInformation.PersonId, c.sessionInformation.PersonType, startDate, endDate, validateSession)
}

func (c *Client) GetTimetableForRange(id int, timeTableType int, startDate time.Time, endDate time.Time, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(id, timeTableType, startDate, endDate, validateSession)
}

func (c *Client) GetOwnClassTimetableForToday(validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(c.sessionInformation.ClassId, 1, time.Time{}, time.Time{}, validateSession)
}

func (c *Client) getOwnClassTimetableFor(date time.Time, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(c.sessionInformation.ClassId, 1, date, date, validateSession)
}

func (c *Client) GetOwnClassTimetableForRange(startDate time.Time, endDate time.Time, validateSession bool) ([]GenericLesson, error) {
	return c.requestTimeTable(c.sessionInformation.ClassId, 1, startDate, endDate, validateSession)
}

func (c *Client) GetHomeworksFor(rangeStart time.Time, rangeEnd time.Time, validateSession bool) ([]Homework, error) {
	if validateSession {
		if err := c.validateSession(); err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.R().SetHeader(
		"Cookie", c.GetCookie(),
	).SetQueryParam(
		"startDate", GetDateUntisFormat(rangeStart),
	).SetQueryParam(
		"endDate", GetDateUntisFormat(rangeEnd),
	).Get(c.BaseUrl + "/WebUntis/api/homeworks/lessons")

	result := gjson.GetBytes(resp.Body(), "data.homeworks")
	if !result.Exists() {
		return nil, eris.Wrap(err, "request didn't return any result")
	}

	// Embed lesson names into homeworks
	var homeworks []Homework

	if err := json.Unmarshal([]byte(result.String()), &homeworks); err != nil {
		return nil, eris.Wrap(err, "homework format incorrect")
	}

	lessonResult := gjson.GetBytes(resp.Body(), "data.lessons")

	if !lessonResult.Exists() {
		return nil, eris.Wrap(err, "request didn't return any result")
	}

	var lessons []LessonWithSubj
	if err := json.Unmarshal([]byte(lessonResult.String()), &lessons); err != nil {
		return nil, eris.Wrap(err, "Lesson format incorrect")
	}

	lessonMap, err := GetLessonMap(lessons)
	if err != nil {
		return nil, eris.Wrap(err, "Error getting lesson from request")
	}

	for i, homework := range homeworks {
		homeworks[i].LessonName = lessonMap[homework.LessonId]
	}

	return homeworks, err
}

func (c *Client) GetSubjects(validateSession bool) ([]Subject, error) {
	resp, err := c.request("getSubjects", "", validateSession)
	if err != nil {
		return nil, eris.Wrap(err, "Error getting subjects")
	}

	result := gjson.GetBytes(resp, "result")
	if !result.Exists() {
		return nil, eris.Wrap(err, "request didn't return any result")
	}

	var subjects []Subject
	if err := json.Unmarshal([]byte(result.String()), &subjects); err != nil {
		return nil, eris.Wrap(err, "Subject format incorrect")
	}

	return subjects, nil
}

func (c *Client) GetTimegrid(validateSession bool) ([]TimeGridLesson, error) {
	resp, err := c.request("getTimegridUnits", "", validateSession)
	if err != nil {
		return nil, err
	}
	var grid []TimeGridLesson

	res := gjson.GetBytes(resp, "result")
	if !res.Exists() {
		log.Println(string(resp))
		return nil, errors.New("key results doesn't exist in answer")
	}

	if err := json.Unmarshal([]byte(res.String()), &grid); err != nil {
		return nil, err
	}
	return grid, nil
}

func (c *Client) GetHomeWorkAndLessons(rangeStart time.Time, rangeEnd time.Time, validateSession bool) ([]byte, error) {
	if validateSession {
		if err := c.validateSession(); err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.R().SetHeader(
		"Cookie", c.GetCookie(),
	).SetQueryParam(
		"startDate", GetDateUntisFormat(rangeStart),
	).SetQueryParam("endDate", GetDateUntisFormat(rangeEnd)).Get(
		c.BaseUrl + "/WebUntis/api/homeworks/lessons",
	)
	if err != nil {
		return nil, eris.Wrap(err, "error getting homeworks and lessons")
	}
	if resp.IsError() {
		return nil, eris.New("server response non 200e")
	}

	return resp.Body(), nil
}

// GetSchoolyears gets all WebUntis Schoolyears
func (c *Client) GetSchoolyears(validateSession bool) ([]SchoolYear, error) {
	data, err := c.request("getSchoolyears", "{}", validateSession)
	if err != nil {
		return nil, err
	}
	resultsJson := gjson.GetBytes(data, "result")
	if !resultsJson.Exists() {
		log.Println(string(data))
		return nil, errors.New("key results doesn't exist in answer")
	}

	var schoolYears []SchoolYear

	err = json.Unmarshal([]byte(resultsJson.String()), &schoolYears)
	if err != nil {
		return nil, err
	}

	// Sort schoolYears by startDate
	sort.Slice(schoolYears, func(i, j int) bool {
		return schoolYears[i].StartDate.Before(schoolYears[j].StartDate)
	})

	return schoolYears, nil
}

// GetLatestSchoolyear gets the latest WebUntis Schoolyear
func (c *Client) GetLatestSchoolyear(validateSession bool) (SchoolYear, error) {
	schoolYears, err := c.GetSchoolyears(validateSession)
	if err != nil {
		return SchoolYear{}, err
	}
	return schoolYears[len(schoolYears)-1], nil
}

func (c *Client) GetClasses(validateSession bool) ([]Class, error) {
	SchoolYear, err := c.GetLatestSchoolyear(validateSession)
	if err != nil {
		return nil, err
	}

	requestData := map[string]int{
		"schoolyearId": SchoolYear.Id,
	}
	respData, err := c.request("getKlassen", ToJsonStr(requestData), validateSession)
	if err != nil {
		return nil, err
	}

	res := gjson.Get(string(respData), "result")
	if !res.Exists() {
		log.Println(string(respData))
		return nil, errors.New("key results doesn't exist in answer")
	}
	var classes []Class

	err = json.Unmarshal([]byte(res.String()), &classes)
	if err != nil {
		return nil, err
	}

	return classes, nil
}

// GetLatestImportTime gets the time when WebUntis last changed its data
func (c *Client) GetLatestImportTime(validateSession bool) (time.Time, error) {
	data, err := c.request("getLatestImportTime", "{}", validateSession)
	if err != nil {
		return time.Time{}, err
	}

	timeInt := gjson.Get(string(data), "result")
	if !timeInt.Exists() {
		return time.Time{}, errors.New("key results doesn't exist in answer")
	}

	return time.Unix(0, timeInt.Int()*int64(time.Millisecond)), nil
}

// GetAbsentLessons returns all the lessons where you were absent including the excused one
func (c *Client) GetAbsentLessons(rangeStart time.Time, rangeEnd time.Time, excuseStateId int, validateSession bool) (Absences, error) {
	if validateSession {
		if err := c.validateSession(); err != nil {
			return Absences{}, err
		}
	}

	resp, err := c.httpClient.R().SetQueryParams(
		map[string]string{
			"startDate":      GetDateUntisFormat(rangeStart),
			"endDate":        GetDateUntisFormat(rangeEnd),
			"studentId":      strconv.Itoa(c.sessionInformation.PersonId),
			"excuseStatusId": strconv.Itoa(excuseStateId),
		},
	).SetHeader(
		"Cookie", c.GetCookie(),
	).Get(
		c.BaseUrl + "/WebUntis/api/classreg/absences/students",
	)

	if err != nil {
		return Absences{}, err
	}

	if resp.IsError() {
		return Absences{}, errors.New("server response non 200")
	}
	if resp.String() == "" {
		return Absences{}, errors.New("server response empty")
	}

	var absences Absences

	res := gjson.Get(resp.String(), "data")
	if !res.Exists() {
		return Absences{}, errors.New("key data doesn't exist")
	}

	if err := json.Unmarshal([]byte(res.String()), &absences); err != nil {
		return Absences{}, err
	}

	return absences, nil
}

func (c *Client) getAccessToken() error {
	generationTime := time.Now()
	otp, err := totp.GenerateCode(c.Secret, generationTime)
	if err != nil {
		log.Println("Error generating otp using token, check secret, %w", err)
		return err
	}

	data := map[string]any{
		"id":     c.Identity,
		"method": "getUserData2017",
		"params": []map[string]any{
			{
				"auth": map[string]any{
					"clientTime": generationTime.UnixMilli(),
					"user":       c.Username,
					"otp":        otp,
				},
			},
		},
		"jsonrpc": "2.0",
	}

	bodyJson, err := json.Marshal(data)
	if err != nil {
		log.Println("error generating json from request data, %w", err)
		return err
	}

	resp, err := c.httpClient.R().SetBody(
		bytes.NewReader(bodyJson),
	).SetHeaders(
		map[string]string{
			"Accept":           "application/json, text/plain, */*",
			"Content-Type":     "application/json",
			"Cache-Control":    "no-cache",
			"Pragma":           "no-cache",
			"X-Requested-With": "XMLHttpRequest",
			"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
			"Host":             "antiope.webuntis.com",
		},
	).SetContentLength(true).Post(
		fmt.Sprintf("%s/WebUntis/jsonrpc_intern.do?m=getUserData2017&school=%s&v=i2.2", c.BaseUrl, c.School),
	)

	if err != nil {
		log.Println("error fetching token, %w", err)
		return err
	}

	if resp.IsError() {
		log.Printf("Error getting token from server, request-body: %v", resp.Body())
		return errors.New("server response non 200")
	}

	c.extractCookieInformation(resp.Header().Get("set-cookie"))
	return nil
}

func (c *Client) extractCookieInformation(cookies string) {
	parts := strings.Split(cookies, ";")
	for _, cookie := range parts {
		cookie = strings.TrimSpace(cookie)
		cookie = strings.Replace(cookie, ";", "", 1)
		keyValue := strings.Split(cookie, "=")
		if len(keyValue) != 2 {
			continue
		}
		if keyValue[0] == "JSESSIONID" {
			c.sessionInformation = SessionInformation{SessionId: keyValue[1]}
			break
		}
	}
}

func (c *Client) GetCookie() string {
	return fmt.Sprintf("schoolname=\"%s\"; JSESSIONID=%s;", "_"+ToBase64(c.School), c.sessionInformation.SessionId)
}

func (c *Client) getHeaders() map[string]string {
	return map[string]string{
		"Cookie":           c.GetCookie(),
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.79 Safari/537.36",
		"Cache-Control":    "no-cache",
		"Pragma":           "no-cache",
		"X-Requested-With": "XMLH<ttpRequest",
	}
}
