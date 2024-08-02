package newt

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	inputFmtStr = map[string]string{
		"week":           `<input type="week" id=%q name=%q value="{{%s}}" %s>`,
		"time":           `<input type="time" id=%q name=%q value="{{%s}}" %s>`,
		"text":           `<input type="text" id=%q name=%q value="{{%s}}" %s>`,
		"search":         `<input type="search" id=%q name=%q value="{{%s}}" %s>`,
		"submit":         `<input type="submit" id=%q name=%q value="{{%s}}" %s>`,
		"reset":          `<input type="reset" id=%q name=%q value="{{%s}}" %s>`,
		"range":          `<input type="range" id=%q name=%q value="{{%s}}" %s>`,
		"radio":          `<input type="radio" id=%q name=%q value="{{%s}}" %s>`,
		"password":       `<input type="password" id=%q name=%q value="{{%s}}" %s>`,
		"number":         `<input type="number" id=%q name=%q value="{{%s}}" %s>`,
		"month":          `<input type="month" id=%q name=%q value="{{%s}}" %s>`,
		"image":          `<input type="image" id=%q name=%q value="{{%s}}" %s>`,
		"hidden":         `<input type="hidden" id=%q name=%q value="{{%s}}" %s>`,
		"file":           `<input type="file" id=%q name=%q value="{{%s}}" %s>`,
		"datetime-local": `<input type="datetime-local" id=%q name=%q value="{{%s}}" %s>`,
		"color":          `<input type="color" id=%q name=%q value="{{%s}}" %s>`,
		"date":           `<input type="date" id=%q name=%q value="{{%s}}" %s>`,
		"url":            `<input type="url" id=%q name=%q value="{{%s}}" %s>`,
		"email":          `<input type="email" id=%q name=%q value="{{%s}}" %s>`,
		"button":         `<input type="button" id=%q name=%q value="{{%s}}" %s>`,
		// Alias of orcid example
		"orcid": `<input type="text" extended-type="orcid" id=%q name=%q value="{{%s}}" pattern="[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9A-Z]">`,
	}
)

// genElementAttrsString will generate a string that contains all the rendered element attributes
// e.g. id, class, title, placeholder text
func genElementAttrString(attributes map[string]string, excludeList []string) string {
	parts := []string{}
	for k, v := range attributes {
		if !inList(k, excludeList) {
			parts = append(parts, fmt.Sprintf("%s=%q", k, v))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// displayElemGen generates a specific Mustache template for displaying
// and element.
func displayElemGen(elem *Element) string {
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
	excludeList := []string{"label", "placeholder"}
	attrs := genElementAttrString(elem.Attributes, excludeList)
	switch elem.Type {
	case "phone":
		return fmt.Sprintf(`<a href="tel:{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	case "url":
		return fmt.Sprintf(`<a href="{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	case "email":
		return fmt.Sprintf(`<a href="mailto:{{%s}}" %s>{{%s}}</a>`, elem.Id, attrs, elem.Id)
	default:
		return fmt.Sprintf(`<span %s>{{%s}}</span>`, attrs, elem.Id)
	}
}

// inputElemGen generates a specific Mustache template for input of an element
func inputElemGen(elem *Element) string {
	var (
		input string
	)
	excludeList := []string{"label"}
	attrs := genElementAttrString(elem.Attributes, excludeList)
	if fmtStr, ok := inputFmtStr[elem.Type]; ok {
		input = fmt.Sprintf(fmtStr, elem.Id, elem.Id, elem.Id, attrs)
	} else {
		input = fmt.Sprintf(`<input type=%q id=%s name=%s value="{{%s}}" %s>`, elem.Type, elem.Id, elem.Id, elem.Id, attrs)
	}
	if label, ok := elem.Attributes["label"]; ok {
		return fmt.Sprintf(`<div><label for=%q>%s</label> %s</div>`, elem.Id, label, input)
	}
	return input
}

// tmplGenCreateForm generations a Mustache partial for a create object web form
func tmplGenCreateForm(out io.Writer, model *Model) error {
	formURL := fmt.Sprintf("/create_%s", model.Id)
	fmt.Fprintf(out, "<form method=%q action=%q>", http.MethodPost, formURL)
	// Build a webform partial
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := inputElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	fmt.Fprintf(out, `<input type="submit" value="create"> <input type="reset" value="reset">`)
	fmt.Fprintf(out, "</form>\n")
	return nil
}

// tmplGenCreateResponse generations a Mustache partial for a create object web form
func tmplGenCreateResponse(out io.Writer, model *Model) error {
	// Build a response partial
	return tmplGenRead(out, model)
}

// tmplGenRead generations a Mustache partial for a read object display element
func tmplGenRead(out io.Writer, model *Model) error {
	// Build a display partial
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := displayElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	return nil
}

// tmplGenUpdateForm generations a Mustache partial for a update object web form
func tmplGenUpdateForm(out io.Writer, model *Model) error {
	// Build a webform partial
	formURL := fmt.Sprintf("/%s_update", model.Id)
	fmt.Fprintf(out, "<form method=%q action=%q>", http.MethodPost, formURL)
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := inputElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	fmt.Fprintf(out, `<input type="submit" value="update"> <input type="reset" value="reset">`)
	fmt.Fprintf(out, "</form>\n")
	return nil
}

// tmplGenUpdateResponse generations a Mustache partial for a update object web form
func tmplGenUpdateResponse(out io.Writer, model *Model) error {
	// Build a response partial
	return tmplGenRead(out, model)
}

// tmplGenDeleteForm generations a Mustache partial for a delete object web form
func tmplGenDeleteForm(out io.Writer, model *Model) error {
	//FIXME: what do we need to delete a record? Just the record id? other fields?
	// Build a webform partial
	formURL := fmt.Sprintf("/%s_delete", model.Id)
	fmt.Fprintf(out, "<form method=%q action=%q>", http.MethodPost, formURL)
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			elem.Attributes["disabled"] = "true"
			if s := inputElemGen(elem); s != "" {
				fmt.Fprintf(out, "\t%s\n", s)
			}
		}
	}
	fmt.Fprintf(out, `<input type="submit" value="delete">`)
	fmt.Fprintf(out, "</form>\n")
	return nil
}

// tmplGenDeleteResponse generations a Mustache partial for a delete object web form
func tmplGenDeleteResponse(out io.Writer, model *Model) error {
	//FIXME: what do we need to delete a record? Just the record id? other fields?
	// Build a webform partial
	fmt.Fprintf(out, "<b>object deleted goes here ...</b>\n")
	return nil
}

// tmplGenList generations a Mustache partial to lists objects
func tmplGenList(out io.Writer, model *Model) error {
	fmt.Fprintf(out, "<ul>\n")
	fmt.Fprintf(out, "{{#%s}}\n", model.Id)
	fmt.Fprintf(out, "\t<li>")
	for i, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := displayElemGen(elem); s != "" {
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

func TmplHeadPartial(out io.Writer, defaultTitle string, cssPath string) error {
	fmt.Fprintf(out, `<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    {{#if page_title}}<title>{{page_title}}</title>{{else}}<title>%s</title>{{/if}}
    <link rel="stylesheet" href="%s">
</head>
`, defaultTitle, cssPath)
	return nil
}

func TmplHeaderPartial(out io.Writer, defaultHeaderLink string, defaultHeaderLinkTitle string, defaultLogoLink string, defaultLogoTitle string) error {
	fmt.Fprintf(out, `<header>
<a href="%s" title="%s"><img src="%s" alt="%s"></a>
</header>
`, defaultHeaderLink, defaultHeaderLinkTitle, defaultLogoLink, defaultLogoTitle)
	return nil
}

func TmplNavPartial(out io.Writer, navElement string) error {
	if navElement != "" {
		fmt.Fprintln(out, navElement)
	} else {
		fmt.Fprintln(out, `<!-- NAV ELEMENT GOES HERE -->`)
	}
	return nil
}

func TmplFooterPartial(out io.Writer, copyrightYear string, copyrightLink string, copyrightText string, licenseLink string, licenseText string, contactAddress string, contactEMail string, contactPhone string) error {
	fmt.Fprintf(out, `<footer>
<span id="copyright">&copy; %s <a href="%s">%s</a></span>
`, copyrightYear, copyrightLink, copyrightText)
	if licenseText != "" {
		fmt.Fprintf(out, `<span id="license"><a href="%s">%s</a></span>
`, licenseLink, licenseText)
	}
	if contactAddress != "" {
		fmt.Fprintf(out, `
<address>%s</address>
`, contactAddress)
	}
	if contactEMail != "" {
		fmt.Fprintf(out, `<span><a href="mailto:%s">Email Us</a></span>
`, contactEMail)
	}
	if contactPhone != "" {
		fmt.Fprintf(out, `<span>Phone: <a href="tel:%s">%s</a></span>
`, contactPhone, contactPhone)
	}
	fmt.Fprintf(out, `</footer>
`)
	return nil
}

type BodyGen func(io.Writer, *Model) error

// tmplPage takes an output buffer and a PartialGen with a function
// signature `func(io.Writer,*Model) error` and renders a webpage
// using the passed in func.
func tmplPage(out io.Writer, model *Model, fn BodyGen) error {
	fmt.Fprintf(out, `<!DOCTYPE html>
<html lang="en-us">
{{>head}}
	<body>
{{>header}}
{{>nav}}
<section>
`)
	err := fn(out, model)
	fmt.Fprintf(out, `
</section>
{{>footer}}
	</body>
</html>`)
	return err
}

// TmplGen takes an io.Writer, model and an action string rendering
// the contents of the model as a Newt handlebars template for the
// provided action. It returns an error value when something
// goes wrong.
func TmplGen(out io.Writer, model *Model, action string) error {
	if model == nil || model.Id == "" {
		return fmt.Errorf("model appears incomplete, aborting")
	}
	var (
		err error
	)
	switch action {
	case "create_form":
		err = tmplPage(out, model, tmplGenCreateForm)
	case "create_response":
		err = tmplPage(out, model, tmplGenCreateResponse)
	case "read":
		err = tmplPage(out, model, tmplGenRead)
	case "update_form":
		err = tmplPage(out, model, tmplGenUpdateForm)
	case "update_response":
		err = tmplPage(out, model, tmplGenUpdateResponse)
	case "delete_form":
		err = tmplPage(out, model, tmplGenDeleteForm)
	case "delete_response":
		err = tmplPage(out, model, tmplGenDeleteResponse)
	case "list":
		err = tmplPage(out, model, tmplGenList)
	default:
		return fmt.Errorf("%q generation is not supported", action)
	}
	return err
}
