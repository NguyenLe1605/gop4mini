// Helper class to load p4info files, get value from p4 protobuf object and create table entry. 
package controller

import (
	"fmt"
	"io/ioutil"

	"github.com/NguyenLe1605/gop4mini/pkg/utils"
	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
	"google.golang.org/protobuf/encoding/prototext"
)

type P4InfoHelper struct {
	P4Info *p4info.P4Info
}

func NewP4InfoHelper(p4InfoFilePath string) (*P4InfoHelper, error) {
	data, err := ioutil.ReadFile(p4InfoFilePath)
	if err != nil {
		return nil, err
	}

	info := &p4info.P4Info{}
	err = prototext.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	return &P4InfoHelper{
		P4Info: info,
	}, nil
}

func (p *P4InfoHelper) GetTableId(tableName string) int32 {
	for _, table := range p.P4Info.Tables {
		if table.Preamble.Name == tableName {
			return int32(table.Preamble.Id)
		}
	}
	return -1
}

func (p *P4InfoHelper) GetActionId(actionName string) int32 {
	for _, action := range p.P4Info.Actions {
		if action.Preamble.Name == actionName {
			return int32(action.Preamble.Id)
		}
	}
	return -1
}

func (p *P4InfoHelper) GetMatchFieldByName(tableName string, name string) (*p4info.MatchField, error) {
	for _, table := range p.P4Info.Tables {
		pre := table.Preamble
		if pre.Name == tableName {
			for _, mf := range table.MatchFields {
				if mf.Name == name {
					return mf, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("%s has no attribute %s", tableName, name)
}

func (p *P4InfoHelper) GetMatchFieldPb(
	tableName string,
	matchFieldName string,
	value any,
) (*p4api.FieldMatch, error) {
	p4InfoMatch, err := p.GetMatchFieldByName(tableName, matchFieldName)
	if err != nil {
		return nil, err
	}

	bitwidth := p4InfoMatch.Bitwidth

	match := &p4api.FieldMatch{
		FieldId: p4InfoMatch.Id,
	}

	switch p4InfoMatch.GetMatchType() {
	case p4info.MatchField_EXACT:
		val, err := utils.Encode(value, int(bitwidth))
		if err != nil {
			return nil, err
		}
		exact := &p4api.FieldMatch_Exact{
			Value: val,
		}
		match.FieldMatchType = &p4api.FieldMatch_Exact_{Exact: exact}
	case p4info.MatchField_LPM:
		tup, _ := value.(struct {
			string
			int32
		})
		val, err := utils.Encode(tup.string, int(bitwidth))
		if err != nil {
			return nil, err
		}
		prefixLen := tup.int32
		lpm := &p4api.FieldMatch_LPM{
			Value:     val,
			PrefixLen: prefixLen,
		}
		match.FieldMatchType = &p4api.FieldMatch_Lpm{Lpm: lpm}
	default:
		return nil, fmt.Errorf("match field not support yet %v", p4InfoMatch.GetMatchType())
	}
	return match, nil
}

func (p *P4InfoHelper) GetActionParam(actionName string, paramName string) *p4info.Action_Param {
	for _, action := range p.P4Info.Actions {
		pre := action.Preamble
		if pre.Name == actionName {
			for _, param := range action.Params {
				if param.Name == paramName {
					return param
				}
			}
		}
	}
	return nil
}

func (p *P4InfoHelper) GetActionParamPb(
	actionName string,
	paramName string,
	value any,
) (*p4api.Action_Param, error) {
	p4InfoParam := p.GetActionParam(actionName, paramName)
	paramId := p4InfoParam.Id
	actionParam := &p4api.Action_Param{
		ParamId: paramId,
	}

	bitwidth := p4InfoParam.Bitwidth

	val, err := utils.Encode(value, int(bitwidth))
	if err != nil {
		return nil, err
	}

	actionParam.Value = val
	return actionParam, nil
}

func (p *P4InfoHelper) BuildTableEntry(
	tableName string,
	matchFields map[string]interface{},
	actionName string,
	actionParams map[string]interface{},
) (*p4api.TableEntry, error) {
	tableId := p.GetTableId(tableName)
	matches := make([]*p4api.FieldMatch, 0)
	for matchFieldName, value := range matchFields {
		field, _ := p.GetMatchFieldPb(tableName, matchFieldName, value)
		matches = append(matches, field)
	}

	params := make([]*p4api.Action_Param, 0)
	for fieldName, value := range actionParams {
		param, err := p.GetActionParamPb(actionName, fieldName, value)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	actionId := p.GetActionId(actionName)
	action := &p4api.Action{
		ActionId: uint32(actionId),
		Params:   params,
	}
	tableEntry := &p4api.TableEntry{
		TableId: uint32(tableId),
		Match:   matches,
		Action: &p4api.TableAction{
			Type: &p4api.TableAction_Action{
				Action: action,
			},
		},
	}
	return tableEntry, nil
}
