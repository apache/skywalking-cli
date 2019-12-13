// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package schema

import (
	"fmt"
	"io"
	"strconv"
)

type AlarmMessage struct {
	StartTime int64  `json:"startTime"`
	Scope     *Scope `json:"scope"`
	ID        string `json:"id"`
	Message   string `json:"message"`
}

type AlarmTrend struct {
	NumOfAlarm []*int `json:"numOfAlarm"`
}

type Alarms struct {
	Msgs  []*AlarmMessage `json:"msgs"`
	Total int             `json:"total"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BasicTrace struct {
	SegmentID     string   `json:"segmentId"`
	EndpointNames []string `json:"endpointNames"`
	Duration      int      `json:"duration"`
	Start         string   `json:"start"`
	IsError       *bool    `json:"isError"`
	TraceIds      []string `json:"traceIds"`
}

type BatchMetricConditions struct {
	Name string   `json:"name"`
	Ids  []string `json:"ids"`
}

type Call struct {
	Source           string        `json:"source"`
	SourceComponents []string      `json:"sourceComponents"`
	Target           string        `json:"target"`
	TargetComponents []string      `json:"targetComponents"`
	ID               string        `json:"id"`
	DetectPoints     []DetectPoint `json:"detectPoints"`
}

type ClusterBrief struct {
	NumOfService  int `json:"numOfService"`
	NumOfEndpoint int `json:"numOfEndpoint"`
	NumOfDatabase int `json:"numOfDatabase"`
	NumOfCache    int `json:"numOfCache"`
	NumOfMq       int `json:"numOfMQ"`
}

type Database struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Step  Step   `json:"step"`
}

type Endpoint struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EndpointInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
}

type IntValues struct {
	Values []*KVInt `json:"values"`
}

type KVInt struct {
	ID    string `json:"id"`
	Value int64  `json:"value"`
}

type KeyValue struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

type Log struct {
	ServiceName         *string     `json:"serviceName"`
	ServiceID           *string     `json:"serviceId"`
	ServiceInstanceName *string     `json:"serviceInstanceName"`
	ServiceInstanceID   *string     `json:"serviceInstanceId"`
	EndpointName        *string     `json:"endpointName"`
	EndpointID          *string     `json:"endpointId"`
	TraceID             *string     `json:"traceId"`
	Timestamp           string      `json:"timestamp"`
	IsError             *bool       `json:"isError"`
	StatusCode          *string     `json:"statusCode"`
	ContentType         ContentType `json:"contentType"`
	Content             *string     `json:"content"`
}

type LogEntity struct {
	Time int64       `json:"time"`
	Data []*KeyValue `json:"data"`
}

type LogQueryCondition struct {
	MetricName        *string     `json:"metricName"`
	ServiceID         *string     `json:"serviceId"`
	ServiceInstanceID *string     `json:"serviceInstanceId"`
	EndpointID        *string     `json:"endpointId"`
	TraceID           *string     `json:"traceId"`
	QueryDuration     *Duration   `json:"queryDuration"`
	State             LogState    `json:"state"`
	StateCode         *string     `json:"stateCode"`
	Paging            *Pagination `json:"paging"`
}

type Logs struct {
	Logs  []*Log `json:"logs"`
	Total int    `json:"total"`
}

type MetricCondition struct {
	Name string  `json:"name"`
	ID   *string `json:"id"`
}

type Node struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   *string `json:"type"`
	IsReal bool    `json:"isReal"`
}

type Pagination struct {
	PageNum   *int  `json:"pageNum"`
	PageSize  int   `json:"pageSize"`
	NeedTotal *bool `json:"needTotal"`
}

type Ref struct {
	TraceID         string  `json:"traceId"`
	ParentSegmentID string  `json:"parentSegmentId"`
	ParentSpanID    int     `json:"parentSpanId"`
	Type            RefType `json:"type"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ServiceInstance struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Attributes   []*Attribute `json:"attributes"`
	Language     Language     `json:"language"`
	InstanceUUID string       `json:"instanceUUID"`
}

type Span struct {
	TraceID      string       `json:"traceId"`
	SegmentID    string       `json:"segmentId"`
	SpanID       int          `json:"spanId"`
	ParentSpanID int          `json:"parentSpanId"`
	Refs         []*Ref       `json:"refs"`
	ServiceCode  string       `json:"serviceCode"`
	StartTime    int64        `json:"startTime"`
	EndTime      int64        `json:"endTime"`
	EndpointName *string      `json:"endpointName"`
	Type         string       `json:"type"`
	Peer         *string      `json:"peer"`
	Component    *string      `json:"component"`
	IsError      *bool        `json:"isError"`
	Layer        *string      `json:"layer"`
	Tags         []*KeyValue  `json:"tags"`
	Logs         []*LogEntity `json:"logs"`
}

type Thermodynamic struct {
	Nodes     [][]*int `json:"nodes"`
	AxisYStep int      `json:"axisYStep"`
}

type TimeInfo struct {
	Timezone         *string `json:"timezone"`
	CurrentTimestamp *int64  `json:"currentTimestamp"`
}

type TopNEntity struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Value int64  `json:"value"`
}

type TopNRecord struct {
	Statement *string `json:"statement"`
	Latency   int64   `json:"latency"`
	TraceID   *string `json:"traceId"`
}

type TopNRecordsCondition struct {
	ServiceID  string    `json:"serviceId"`
	MetricName string    `json:"metricName"`
	TopN       int       `json:"topN"`
	Order      Order     `json:"order"`
	Duration   *Duration `json:"duration"`
}

type Topology struct {
	Nodes []*Node `json:"nodes"`
	Calls []*Call `json:"calls"`
}

type Trace struct {
	Spans []*Span `json:"spans"`
}

type TraceBrief struct {
	Data  []*BasicTrace `json:"data"`
	Total int           `json:"total"`
}

type TraceQueryCondition struct {
	ServiceID         *string     `json:"serviceId"`
	ServiceInstanceID *string     `json:"serviceInstanceId"`
	TraceID           *string     `json:"traceId"`
	EndpointID        *string     `json:"endpointId"`
	EndpointName      *string     `json:"endpointName"`
	QueryDuration     *Duration   `json:"queryDuration"`
	MinTraceDuration  *int        `json:"minTraceDuration"`
	MaxTraceDuration  *int        `json:"maxTraceDuration"`
	TraceState        TraceState  `json:"traceState"`
	QueryOrder        QueryOrder  `json:"queryOrder"`
	Paging            *Pagination `json:"paging"`
}

type ContentType string

const (
	ContentTypeText ContentType = "TEXT"
	ContentTypeJSON ContentType = "JSON"
	ContentTypeNone ContentType = "NONE"
)

var AllContentType = []ContentType{
	ContentTypeText,
	ContentTypeJSON,
	ContentTypeNone,
}

func (e ContentType) IsValid() bool {
	switch e {
	case ContentTypeText, ContentTypeJSON, ContentTypeNone:
		return true
	}
	return false
}

func (e ContentType) String() string {
	return string(e)
}

func (e *ContentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContentType", str)
	}
	return nil
}

func (e ContentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DetectPoint string

const (
	DetectPointClient DetectPoint = "CLIENT"
	DetectPointServer DetectPoint = "SERVER"
	DetectPointProxy  DetectPoint = "PROXY"
)

var AllDetectPoint = []DetectPoint{
	DetectPointClient,
	DetectPointServer,
	DetectPointProxy,
}

func (e DetectPoint) IsValid() bool {
	switch e {
	case DetectPointClient, DetectPointServer, DetectPointProxy:
		return true
	}
	return false
}

func (e DetectPoint) String() string {
	return string(e)
}

func (e *DetectPoint) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DetectPoint(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DetectPoint", str)
	}
	return nil
}

func (e DetectPoint) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Language string

const (
	LanguageUnknown Language = "UNKNOWN"
	LanguageJava    Language = "JAVA"
	LanguageDotnet  Language = "DOTNET"
	LanguageNodejs  Language = "NODEJS"
	LanguagePython  Language = "PYTHON"
	LanguageRuby    Language = "RUBY"
)

var AllLanguage = []Language{
	LanguageUnknown,
	LanguageJava,
	LanguageDotnet,
	LanguageNodejs,
	LanguagePython,
	LanguageRuby,
}

func (e Language) IsValid() bool {
	switch e {
	case LanguageUnknown, LanguageJava, LanguageDotnet, LanguageNodejs, LanguagePython, LanguageRuby:
		return true
	}
	return false
}

func (e Language) String() string {
	return string(e)
}

func (e *Language) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Language(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Language", str)
	}
	return nil
}

func (e Language) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type LogState string

const (
	LogStateAll     LogState = "ALL"
	LogStateSuccess LogState = "SUCCESS"
	LogStateError   LogState = "ERROR"
)

var AllLogState = []LogState{
	LogStateAll,
	LogStateSuccess,
	LogStateError,
}

func (e LogState) IsValid() bool {
	switch e {
	case LogStateAll, LogStateSuccess, LogStateError:
		return true
	}
	return false
}

func (e LogState) String() string {
	return string(e)
}

func (e *LogState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LogState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LogState", str)
	}
	return nil
}

func (e LogState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NodeType string

const (
	NodeTypeService  NodeType = "SERVICE"
	NodeTypeEndpoint NodeType = "ENDPOINT"
	NodeTypeUser     NodeType = "USER"
)

var AllNodeType = []NodeType{
	NodeTypeService,
	NodeTypeEndpoint,
	NodeTypeUser,
}

func (e NodeType) IsValid() bool {
	switch e {
	case NodeTypeService, NodeTypeEndpoint, NodeTypeUser:
		return true
	}
	return false
}

func (e NodeType) String() string {
	return string(e)
}

func (e *NodeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NodeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NodeType", str)
	}
	return nil
}

func (e NodeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Order string

const (
	OrderAsc Order = "ASC"
	OrderDes Order = "DES"
)

var AllOrder = []Order{
	OrderAsc,
	OrderDes,
}

func (e Order) IsValid() bool {
	switch e {
	case OrderAsc, OrderDes:
		return true
	}
	return false
}

func (e Order) String() string {
	return string(e)
}

func (e *Order) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Order(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Order", str)
	}
	return nil
}

func (e Order) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type QueryOrder string

const (
	QueryOrderByStartTime QueryOrder = "BY_START_TIME"
	QueryOrderByDuration  QueryOrder = "BY_DURATION"
)

var AllQueryOrder = []QueryOrder{
	QueryOrderByStartTime,
	QueryOrderByDuration,
}

func (e QueryOrder) IsValid() bool {
	switch e {
	case QueryOrderByStartTime, QueryOrderByDuration:
		return true
	}
	return false
}

func (e QueryOrder) String() string {
	return string(e)
}

func (e *QueryOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = QueryOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid QueryOrder", str)
	}
	return nil
}

func (e QueryOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RefType string

const (
	RefTypeCrossProcess RefType = "CROSS_PROCESS"
	RefTypeCrossThread  RefType = "CROSS_THREAD"
)

var AllRefType = []RefType{
	RefTypeCrossProcess,
	RefTypeCrossThread,
}

func (e RefType) IsValid() bool {
	switch e {
	case RefTypeCrossProcess, RefTypeCrossThread:
		return true
	}
	return false
}

func (e RefType) String() string {
	return string(e)
}

func (e *RefType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RefType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RefType", str)
	}
	return nil
}

func (e RefType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Scope string

const (
	ScopeService                 Scope = "Service"
	ScopeServiceInstance         Scope = "ServiceInstance"
	ScopeEndpoint                Scope = "Endpoint"
	ScopeServiceRelation         Scope = "ServiceRelation"
	ScopeServiceInstanceRelation Scope = "ServiceInstanceRelation"
	ScopeEndpointRelation        Scope = "EndpointRelation"
)

var AllScope = []Scope{
	ScopeService,
	ScopeServiceInstance,
	ScopeEndpoint,
	ScopeServiceRelation,
	ScopeServiceInstanceRelation,
	ScopeEndpointRelation,
}

func (e Scope) IsValid() bool {
	switch e {
	case ScopeService,
		ScopeServiceInstance,
		ScopeEndpoint,
		ScopeServiceRelation,
		ScopeServiceInstanceRelation,
		ScopeEndpointRelation:
		return true
	}
	return false
}

func (e Scope) String() string {
	return string(e)
}

func (e *Scope) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Scope(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Scope", str)
	}
	return nil
}

func (e Scope) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Step string

const (
	StepMonth  Step = "MONTH"
	StepDay    Step = "DAY"
	StepHour   Step = "HOUR"
	StepMinute Step = "MINUTE"
	StepSecond Step = "SECOND"
)

var AllStep = []Step{
	StepMonth,
	StepDay,
	StepHour,
	StepMinute,
	StepSecond,
}

func (e Step) IsValid() bool {
	switch e {
	case StepMonth, StepDay, StepHour, StepMinute, StepSecond:
		return true
	}
	return false
}

func (e Step) String() string {
	return string(e)
}

func (e *Step) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Step(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Step", str)
	}
	return nil
}

func (e Step) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TraceState string

const (
	TraceStateAll     TraceState = "ALL"
	TraceStateSuccess TraceState = "SUCCESS"
	TraceStateError   TraceState = "ERROR"
)

var AllTraceState = []TraceState{
	TraceStateAll,
	TraceStateSuccess,
	TraceStateError,
}

func (e TraceState) IsValid() bool {
	switch e {
	case TraceStateAll, TraceStateSuccess, TraceStateError:
		return true
	}
	return false
}

func (e TraceState) String() string {
	return string(e)
}

func (e *TraceState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TraceState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TraceState", str)
	}
	return nil
}

func (e TraceState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
