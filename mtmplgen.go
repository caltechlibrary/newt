package newt

import (
	"fmt"
	"io"
	"strings"
)

var (
	inputFmtStr = map[string]string{
    	"input[type=week]": `<input type="week" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=time]": `<input type="time" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=text]": `<input type="text" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=search]": `<input type="search" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=submit]": `<input type="submit" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=reset]": `<input type="reset" id=%q name=%q value="{{%s}}" %s>`, 
    	"input[type=range]": `<input type="range" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=radio]": `<input type="radio" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=password]": `<input type="password" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=number]": `<input type="number" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=month]": `<input type="month" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=image]": `<input type="image" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=hidden]": `<input type="hidden" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=file]": `<input type="file" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=datetime-local]": `<input type="datetime-local" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=color]": `<input type="color" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=date]": `<input type="date" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=url]": `<input type="url" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=email]": `<input type="email" id=%q name=%q value="{{%s}}" %s>`,
    	"input[type=button]": `<input type="button" id=%q name=%q value="{{%s}}" %s>`,
	}
)


// inList takes a string and compares it with a list of strings. It
// returns true when a match is found, false otherwise
func inList(target string, list []string) bool {
	for _, val := range list {
		if val == target {
			return true
		}
	}
	return false
}

// genElementAttrsString will generate a string that contains all the rendered element attributes
// e.g. id, class, title, placeholder text
func genElementAttrString(attributes map[string]string, excludeList []string) string {
	parts := []string{}
	for k, v := range attributes {
		if ! inList(k, excludeList) {
			parts = append(parts, fmt.Sprintf("%s=%q", k, v))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// MDisplayElemGen generates a specific Mustache template for displaying
// and element.
func MDisplayElemGen(elem *Element) string {
	// Apply a class value based on element id.
	if class, ok := elem.Attributes["class"]; ok {
		if class != "" {
			class += ", " + elem.Id
		} else {
			class = elem.Id
		}
		elem.Attributes["class"] = class
	} else {
		elem.Attributes["class"] = elem.Id
	}
	// Build out our element markup
	excludeList := []string{ "label", "placeholder" } 
	attrs := genElementAttrString(elem.Attributes, excludeList)
	switch elem.Type {
	case "input[type=phone]":
		return fmt.Sprintf(`<a href="tel:{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	case "input[type=url]":
		return fmt.Sprintf(`<a href="{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	case "input[type=email]":
		return fmt.Sprintf(`<a href="mailto:{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	default:
		return fmt.Sprintf(`<span %s>{{%s}}</span>`, attrs, elem.Id)
	}
}

// MInputElemGen generates a specific Mustache template for input of an element
func MInputElemGen(elem *Element) string {
	var (
		input string
	)
	if elem.Type == "button" {
		elem.Type = "input[type=button]"
	}
	if elem.Type == "input" {
		elem.Type = "input[type=text]"
	}
	excludeList := []string{ "label" } 
	attrs := genElementAttrString(elem.Attributes, excludeList)
	if fmtStr, ok := inputFmtStr[elem.Type]; ok {
		input = fmt.Sprintf(fmtStr, elem.Id, elem.Id, elem.Id, attrs)
	} else {
		inputType := strings.TrimSuffix(strings.TrimPrefix(elem.Type, "input[type="), "]")
		input = fmt.Sprintf(`<input type=%q id=%s name=%s value="{{%s}}" %s>`, inputType, elem.Id, elem.Id, elem.Id, attrs)
	}
	if label, ok := elem.Attributes["label"]; ok {
		return fmt.Sprintf(`<div><label for=%q>%s</label> %s</div>`, elem.Id, label, input)
	}
	return input
}

// mTmplGenCreate generations a Mustache partial for a create object web form
func mTmplGenCreate(out io.Writer, model *NewtModel) error {
	// Build a webform partial
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MInputElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	return nil
}

// mTmplGenRead generations a Mustache partial for a read object display element
func mTmplGenRead(out io.Writer, model *NewtModel) error {
	// Build a display partial
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MDisplayElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	return nil
}

// mTmplGenUpdate generations a Mustache partial for a update object web form
func mTmplGenUpdate(out io.Writer, model *NewtModel) error {
	// Build a webform partial
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MInputElemGen(elem); s != "" {
				fmt.Fprint(out,"\t%s\n", s)
			}
		}
	}
	return nil
}

// mTmplGenDelete generations a Mustache partial for a delete object web form
func mTmplGenDelete(out io.Writer, model *NewtModel) error {
	//FIXME: what do we need to delete a record? Just the record id? other fields?
	return fmt.Errorf("mTmplGenDelete(out, %q) not implemented", model.Id)
}

// mTmplGenList generations a Mustache partial to lists objects 
func mTmplGenList(out io.Writer, model *NewtModel) error {
	fmt.Fprintf(out, "<ul>\n")
	fmt.Fprintf(out, "{{#%s}}\n", model.Id)
	fmt.Fprintf(out, "\t<li>")
	for i, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MDisplayElemGen(elem); s != "" {
				if i > 0 {
					fmt.Fprintf(out, " ")
				}
				fmt.Fprintf(out, "%s", s)
			}
		}
	}
	fmt.Fprintf(out, "</li>\n")
	fmt.Fprintf(out, "{{/%s}}\n", model.Id)
	fmt.Fprintf(out, "</ul>\n")
	return nil
}

// MTmplGen takes an io.Writer, an action string and model and renders
// the contents of the model as a Newt Mustache template for the provided
// action. It returns an error value when something goes wrong.
func MTmplGen(out io.Writer, action string, model *NewtModel) error {
	if model == nil || model.Id == "" {
		return fmt.Errorf("model appears incomplete, aborting")
	}
	switch action {
	case "create":
		return mTmplGenCreate(out, model)
	case "read":
		return mTmplGenRead(out, model)
	case "update":
		return mTmplGenUpdate(out, model)
	case "delete":
		return mTmplGenDelete(out, model)
	case "list":
		return mTmplGenList(out, model)
	default:
		return fmt.Errorf("%q generation is not supported", action)
	}
}
