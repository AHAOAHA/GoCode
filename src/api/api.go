/******************************************************************
 * Copyright(C) 2020-2020. All right reserved.
 *
 * Filename: api.go
 * Author: ahaoozhang
 * Date: 2020-03-15 19:27:57 (Sunday)
 * Describe:
 ******************************************************************/
package api

import (
	"GradeManager/src/dao"
	DataCenter "GradeManager/src/proto"
	"errors"
	"strconv"
)

var TeacherCache *map[uint64]DataCenter.TeacherInfo
var StudentCache *map[uint64]DataCenter.StudentInfo
var CollegeCache *map[uint64]DataCenter.CollegeInfo
var MajorCache *map[uint64]DataCenter.MajorInfo

func init() {
	if TeacherCache == nil {
		TeacherCache = new(map[uint64]DataCenter.TeacherInfo)
	}
	if StudentCache == nil {
		StudentCache = new(map[uint64]DataCenter.StudentInfo)
	}
	if CollegeCache == nil {
		CollegeCache = new(map[uint64]DataCenter.CollegeInfo)
	}
	if MajorCache == nil {
		MajorCache = new(map[uint64]DataCenter.MajorInfo)
	}
}

// interface: proto Datacenter.TeacherInfo, every call update TeacherCache.
func GetAllTeacherList() (map[uint64]DataCenter.TeacherInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `teacher`")
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.TeacherInfo)

	for _, v := range sm {
		// teacher_uid college_uid password name sex NRIC status create_time
		var teacher_uid, college_uid uint64
		var sta int
		sta, _ = strconv.Atoi(string(v["status"].([]uint8)))
		crtt, _ := strconv.Atoi(string(v["create_time"].([]uint8)))
		teacher_uid, _ = strconv.ParseUint(string(v["teacher_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		m[teacher_uid] = DataCenter.TeacherInfo{
			TeacherUid: teacher_uid,
			CollegeUid: college_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.TeacherInfo_STATUS(sta),
			CreateTime: uint32(crtt),
		}
	}

	// update cache
	TeacherCache = &m
	return m, nil
}

// interface: proto Datacenter.StudentInfo, every call update StudentCache.
func GetAllStudentList() (map[uint64]DataCenter.StudentInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `student`")
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		// student_uid class_uid college_uid major_uid password name sex NRIC status create_time
		var student_uid, class_uid, college_uid, major_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
		class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}

	// update cache
	StudentCache = &m
	return m, nil
}

// Query student info by teacher_uid, name, college_uid, college_name.
func GetTeacherListByTeacherUid(teacher_uid uint64) (DataCenter.TeacherInfo, error) {
	var rsp DataCenter.TeacherInfo
	var ok bool
	m := make(map[uint64]DataCenter.TeacherInfo)
	if TeacherCache != nil {
		rsp, ok = (*TeacherCache)[teacher_uid]
		if ok {
			return rsp, nil
		}
	}

	// TeacherCache nil or teacher_uid not exist, query from database.
	db_m, err := dao.DataBase.Queryf("select * from `teacher` where `teacher_uid`='%d'", teacher_uid)
	if err != nil || len(db_m) != 1 {
		return rsp, errors.New("query teacher Info err")
	}
	var college_uid uint64
	var sta int
	sta, _ = strconv.Atoi(string(db_m[0]["status"].([]uint8)))
	crtt, _ := strconv.Atoi(string(db_m[0]["create_time"].([]uint8)))
	college_uid, _ = strconv.ParseUint(string(db_m[0]["college_uid"].([]uint8)), 10, 64)
	m[teacher_uid] = DataCenter.TeacherInfo{
		TeacherUid: teacher_uid,
		CollegeUid: college_uid,
		Name:       string(db_m[0]["name"].([]uint8)),
		Password:   string(db_m[0]["password"].([]uint8)),
		Sex:        string(db_m[0]["sex"].([]uint8)),
		NRIC:       string(db_m[0]["NRIC"].([]uint8)),
		Status:     DataCenter.TeacherInfo_STATUS(sta),
		CreateTime: uint32(crtt),
	}

	// update TeacherCache
	if TeacherCache != nil {
		(*TeacherCache)[teacher_uid] = m[teacher_uid]
	}
	rsp = m[teacher_uid]
	return rsp, nil
}

func GetTeacherListByNRIC(NRIC string) (map[uint64]DataCenter.TeacherInfo, error) {
	m := make(map[uint64]DataCenter.TeacherInfo)
	db_m, err := dao.DataBase.Queryf("select * from `teacher` where `NRIC`='%s'", NRIC)
	if err != nil || len(db_m) != 1 {
		return nil, errors.New("query teacher Info err")
	}
	var college_uid, teacher_uid uint64
	var sta int
	sta, _ = strconv.Atoi(string(db_m[0]["status"].([]uint8)))
	crtt, _ := strconv.Atoi(string(db_m[0]["create_time"].([]uint8)))
	college_uid, _ = strconv.ParseUint(string(db_m[0]["college_uid"].([]uint8)), 10, 64)
	teacher_uid, _ = strconv.ParseUint(string(db_m[0]["teacher_uid"].([]uint8)), 10, 64)
	m[teacher_uid] = DataCenter.TeacherInfo{
		TeacherUid: teacher_uid,
		CollegeUid: college_uid,
		Name:       string(db_m[0]["name"].([]uint8)),
		Password:   string(db_m[0]["password"].([]uint8)),
		Sex:        string(db_m[0]["sex"].([]uint8)),
		NRIC:       string(db_m[0]["NRIC"].([]uint8)),
		Status:     DataCenter.TeacherInfo_STATUS(sta),
		CreateTime: uint32(crtt),
	}
	return m, nil
}

// Without Cache.
func GetTeacherListByTeacherName(teacher_name string) (map[uint64]DataCenter.TeacherInfo, error) {
	dbm, err := dao.DataBase.Queryf("select * from `teacher` where `name`='%s'", teacher_name)
	if err != nil || len(dbm) == 0 {
		return nil, errors.New("query teacher info by name err")
	}

	m := make(map[uint64]DataCenter.TeacherInfo)

	for _, v := range dbm {
		var teacher_uid, college_uid uint64
		var sta int
		sta, _ = strconv.Atoi(string(v["status"].([]uint8)))
		crtt, _ := strconv.Atoi(string(v["create_time"].([]uint8)))
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		teacher_uid, err = strconv.ParseUint(string(v["teacher_uid"].([]uint8)), 10, 64)
		m[teacher_uid] = DataCenter.TeacherInfo{
			TeacherUid: teacher_uid,
			CollegeUid: college_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.TeacherInfo_STATUS(sta),
			CreateTime: uint32(crtt),
		}
	}

	return m, nil
}

// Without cache.
func GetTeacherListByCollegeUid(college_uid uint64) (map[uint64]DataCenter.TeacherInfo, error) {
	dbm, err := dao.DataBase.Queryf("select * from `teacher` where `college_uid`='%d'", college_uid)
	if err != nil || len(dbm) == 0 {
		return nil, errors.New("query teacher info by college_uid err")
	}

	m := make(map[uint64]DataCenter.TeacherInfo)

	for _, v := range dbm {
		var teacher_uid, college_uid uint64
		var sta int
		sta, _ = strconv.Atoi(string(v["status"].([]uint8)))
		crtt, _ := strconv.Atoi(string(v["create_time"].([]uint8)))
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		teacher_uid, err = strconv.ParseUint(string(v["teacher_uid"].([]uint8)), 10, 64)
		m[teacher_uid] = DataCenter.TeacherInfo{
			TeacherUid: teacher_uid,
			CollegeUid: college_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.TeacherInfo_STATUS(sta),
			CreateTime: uint32(crtt),
		}
	}

	return m, nil
}

// Without cache.
func GetTeacherListByCollegeName(college_name string) (map[uint64]DataCenter.TeacherInfo, error) {
	cdbm, err := dao.DataBase.Queryf("select * from `college` where `name`='%s'", college_name)
	if err != nil || len(cdbm) != 1 {
		return nil, errors.New("query teacher info by college_uid err")
	}
	college_uid, _ := strconv.ParseUint(string(cdbm[0]["college_uid"].([]uint8)), 10, 64)

	return GetTeacherListByCollegeUid(college_uid)
}

func GetALlCollegeName() ([]string, error) {
	cdbm, err := dao.DataBase.Queryf("select `name` from `college`")
	if err != nil || len(cdbm) == 0 {
		return nil, errors.New("query teacher info by college_uid err")
	}

	var college_name []string
	for _, val := range cdbm {
		college_name = append(college_name, string(val["name"].([]uint8)))
	}

	return college_name, nil
}

func GetNotice() (DataCenter.NoticeInfo, error) {
	var n DataCenter.NoticeInfo
	cdbm, err := dao.DataBase.Queryf("SELECT * from `notice` where `notice_uid` = (SELECT max(`notice_uid`) FROM notice)")
	if err != nil || len(cdbm) != 1 {
		return n, errors.New("query notice err")
	}

	n = DataCenter.NoticeInfo{
		Title:  string(cdbm[0]["title"].([]uint8)),
		Notice: string(cdbm[0]["data"].([]uint8)),
	}
	return n, nil
}

func GetCollegeUidByName(college_name string) (uint64, error) {
	cdbm, err := dao.DataBase.Queryf("SELECT 'college_uid' from `college` where `name`='%s'", college_name)
	if err != nil || len(cdbm) != 1 {
		return 0, err
	}

	college_uid, _ := strconv.ParseUint(string(cdbm[0]["college_uid"].([]uint8)), 10, 64)

	return college_uid, nil
}

func GetAllMajerName() ([]string, error) {
	cdbm, err := dao.DataBase.Queryf("select `name` from `major`")
	if err != nil || len(cdbm) == 0 {
		return nil, errors.New("query major err")
	}

	var major_name []string
	for _, val := range cdbm {
		major_name = append(major_name, string(val["name"].([]uint8)))
	}

	return major_name, nil
}

func GetMajorUidByName(major_name string) (uint64, error) {
	cdbm, err := dao.DataBase.Queryf("SELECT `major_uid` from `major` where `name`='%s'", major_name)
	if err != nil || len(cdbm) != 1 {
		return 0, err
	}

	major_uid, _ := strconv.ParseUint(string(cdbm[0]["major_uid"].([]uint8)), 10, 64)

	return major_uid, nil
}

func GetAllClassName() ([]string, error) {
	cdbm, err := dao.DataBase.Queryf("select `name` from `class`")
	if err != nil || len(cdbm) == 0 {
		return nil, errors.New("query class err")
	}

	var major_name []string
	for _, val := range cdbm {
		major_name = append(major_name, string(val["name"].([]uint8)))
	}

	return major_name, nil
}

func GetStudentByStudentUid(student_uid uint64) (DataCenter.StudentInfo, error) {
	val, ok := (*StudentCache)[student_uid]
	if ok {
		return val, nil
	}

	var s DataCenter.StudentInfo
	sm, err := dao.DataBase.Queryf("select * from `student` where `student_uid`='%d'", student_uid)
	if err != nil || len(sm) != 1 {
		return s, err
	}

	v := sm[0]
	// teacher_uid college_uid password name sex NRIC status create_time
	var college_uid, major_uid, class_uid uint64
	var sta int
	sta, _ = strconv.Atoi(string(v["status"].([]uint8)))
	crtt, _ := strconv.Atoi(string(v["create_time"].([]uint8)))
	college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
	class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
	major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
	s = DataCenter.StudentInfo{
		StudentUid: student_uid,
		CollegeUid: college_uid,
		MajorUid:   major_uid,
		ClassUid:   class_uid,
		Name:       string(v["name"].([]uint8)),
		Password:   string(v["password"].([]uint8)),
		Sex:        string(v["sex"].([]uint8)),
		NRIC:       string(v["NRIC"].([]uint8)),
		Status:     DataCenter.StudentInfo_STATUS(sta),
		CreateTime: int32(crtt),
	}
	return s, nil
}

func GetStudentByNRIC(NRIC string) (map[uint64]DataCenter.StudentInfo, error) {

	sm, err := dao.DataBase.Queryf("select * from `student` where `NRIC`='%s'", NRIC)
	if err != nil {
		return nil, err
	}

	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		// student_uid class_uid college_uid major_uid password name sex NRIC status create_time
		var class_uid, student_uid, college_uid, major_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
		class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}
	return m, nil
}

func GetClassUidByName(class_name string) (uint64, error) {
	m, err := dao.DataBase.Queryf("select `class_uid` from `class` where `name`='%s'", class_name)
	if err != nil || len(m) != 1 {
		return 0, err
	}

	class_uid_str := string(m[0]["class_uid"].([]uint8))
	class_uid, _ := strconv.ParseUint(class_uid_str, 10, 64)
	return class_uid, nil
}

func GetStudentListByClassUid(class_uid uint64) (map[uint64]DataCenter.StudentInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `student` where `class_uid`='%d'", class_uid)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		// student_uid class_uid college_uid major_uid password name sex NRIC status create_time
		var student_uid, college_uid, major_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}

	return m, nil
}

func GetStudentListByMajorUid(major_uid uint64) (map[uint64]DataCenter.StudentInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `student` where `major_uid`='%d'", major_uid)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		// student_uid class_uid college_uid major_uid password name sex NRIC status create_time
		var student_uid, class_uid, college_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}

	return m, nil
}

func GetStudentListByCollegeUid(college_uid uint64) (map[uint64]DataCenter.StudentInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `student` where `college_uid`='%d'", college_uid)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		// student_uid class_uid college_uid major_uid password name sex NRIC status create_time
		var student_uid, class_uid, major_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
		major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       string(v["name"].([]uint8)),
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}

	return m, nil
}

func GetStudentByName(student_name string) (map[uint64]DataCenter.StudentInfo, error) {
	sm, err := dao.DataBase.Queryf("select * from `student` where `name`='%s'", student_name)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]DataCenter.StudentInfo)

	for _, v := range sm {
		var student_uid, class_uid, college_uid, major_uid, crtt uint64
		var sta int
		crtt, _ = strconv.ParseUint(string(v["create_time"].([]uint8)), 10, 64)
		student_uid, _ = strconv.ParseUint(string(v["student_uid"].([]uint8)), 10, 64)
		college_uid, _ = strconv.ParseUint(string(v["college_uid"].([]uint8)), 10, 64)
		major_uid, _ = strconv.ParseUint(string(v["major_uid"].([]uint8)), 10, 64)
		class_uid, _ = strconv.ParseUint(string(v["class_uid"].([]uint8)), 10, 64)
		m[student_uid] = DataCenter.StudentInfo{
			StudentUid: student_uid,
			CollegeUid: college_uid,
			MajorUid:   major_uid,
			ClassUid:   class_uid,
			Name:       student_name,
			Password:   string(v["password"].([]uint8)),
			Sex:        string(v["sex"].([]uint8)),
			NRIC:       string(v["NRIC"].([]uint8)),
			Status:     DataCenter.StudentInfo_STATUS(sta),
			CreateTime: int32(crtt),
		}
	}

	return m, nil
}

func GetNamebyUid(uid uint64, table string, field string) (string, error) {
	m, err := dao.DataBase.Queryf("select `name` from `%s` where `%s`='%d'", table, field, uid)
	if err != nil || len(m) != 1 {
		return "", err
	}

	return string(m[0]["name"].([]uint8)), nil
}