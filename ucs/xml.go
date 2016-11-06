package ucs

import (
    //"fmt"
    "strings"
    "strconv"
    "github.com/beevik/etree"
)

func getValueOfAttribute(array []string, element string) string {
    for _, x := range array {
        if strings.Split(x,"|")[0] == element {
            return strings.Split(x,"|")[1]
        }
    }
    return ""
}

func getPositionOfAttribute(array []string, element string) int {
    for pos, x := range array {
        if strings.Split(x,"|")[0] == element {
            return pos
        }
    }
    return 0
}

func getBulkAttribute(xml string, root string, attributes []string) []string {
    var returns []string
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }

    r := doc.SelectElement(root)
    for _, x := range attributes {
        test := []string{x, r.SelectAttrValue(x, "unknown")}
        returns = append(returns, strings.Join(test[:],"|"))
    }

    return returns
}

func getElementAttributeCumlative(xml string, roots []string, element string, attribute string) string {
    return ""
}

func getElementArray(xml string, roots []string, element string, attribute []string) map[string]string {
    tmp := make(map[string]string)
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }

	var root *etree.Element
	tmpCount := 0
    for _, x := range roots {
		if tmpCount == 0 {
			root = doc.SelectElement(x)
			tmpCount += 1
		} else {
			root = root.SelectElement(x)
		}
	}

    blade := 0
    if root != nil {
        if len(root.SelectElements(element)) > 0 {
            for _,x := range root.SelectElements(element) {
                for _,y := range attribute {
                    tmp["UCS_COMPUTE_"+strconv.Itoa(blade)+"_"+strings.ToUpper(y)] = x.SelectAttrValue(y,"unknown")
                }
                blade += 1
            }
        }
    }
    return tmp
}

func getElementOccuranceVariable(xml string,roots []string,element string,match []string,attribute string) int {
    returns := 0
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }

	var root *etree.Element
	tmpCount := 0
    for _, x := range roots {
		if tmpCount == 0 {
			root = doc.SelectElement(x)
			tmpCount += 1
		} else {
			root = root.SelectElement(x)
		}
	}

    if root != nil {
        if len(root.SelectElements(element)) > 0 {
            for _,x := range root.SelectElements(element) {
                if len(match) == 2 {
                    if x.SelectAttrValue(match[0],"unknown") == match[1] {
                        num,err := strconv.Atoi(x.SelectAttrValue(attribute,"0"))
                        if err == nil {
                            returns += num
                        }
                    }
                }
            }
        }
    }
    return returns
}

func getElementOccurence(xml string,roots []string,element string,attribute string) map[string]int {
    returns := make(map[string]int)
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }
	var root *etree.Element
	tmpCount := 0
    for _, x := range roots {
		if tmpCount == 0 {
			root = doc.SelectElement(x)
			tmpCount += 1
		} else {
			root = root.SelectElement(x)
		}
	}
    if root != nil {
        if len(root.SelectElements(element)) > 0 {
            for _,x := range root.SelectElements(element) {
                returns[x.SelectAttrValue(attribute,"unknown")] += 1
            }
        }
    }
    return returns
}

func getAttribute(xml string, root string, attribute string) string {
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }
    r := doc.SelectElement(root)
    return r.SelectAttrValue(attribute, "unknown")
}

func getElementCount(xml string, roots []string, element string) int {
	count := 0
    doc := etree.NewDocument()
    if err := doc.ReadFromBytes([]byte(xml)); err != nil {
        panic(err)
    }
	var root *etree.Element
	tmpCount := 0
	for _, x := range roots {
		if tmpCount == 0 {
			root = doc.SelectElement(x)
			tmpCount += 1
		} else {
			root = root.SelectElement(x)
		}
	}

	if root != nil {
		return len(root.SelectElements(element))
	}
	return count
}

func xmlReplaceString(xml string, needle string, new string) string {
    return strings.Replace(xml, "#"+needle+"#", new, -1)
}

func xmlReplaceStringArray(xml string, needle []string, new []string) string {
    if len(needle) == len(new) {
        for a := 0; a < len(needle); a++ {
            xml = strings.Replace(xml, "#"+strings.ToUpper(needle[a])+"#", new[a], -1)
        }
    }
    return xml
}
