package grafana

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gapi "github.com/nytm/go-grafana-api"
)

func ResourceAnnotation() *schema.Resource {
	return &schema.Resource{
		Create: CreateAnnotation,
		Read:   ReadAnnotation,
		Update: UpdateAnnotation,
		Delete: DeleteAnnotation,

		Schema: map[string]*schema.Schema{
			"text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func CreateAnnotation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	annotation, err := makeAnnotation(d)
	if err != nil {
		return err
	}

	id, err := client.NewAnnotation(annotation)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(id, 10))

	return ReadAnnotation(d, meta)
}

func ReadAnnotation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid id: %#v", idStr)
	}

	// TODO: handle params & granular arguments
	params := url.Values{}

	as, err := client.Annotations(params)
	if err != nil {
		return err
	}

	annotation := gapi.Annotation{}
	for _, a := range as {
		if a.ID == id {
			annotation = a
			break
		}
	}

	d.Set("id", annotation.ID)
	d.Set("text", annotation.Text)

	return nil
}

func UpdateAnnotation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	annotation, err := makeAnnotation(d)
	if err != nil {
		return err
	}

	_, err = client.UpdateAnnotation(annotation.ID, annotation)

	return err
}

func DeleteAnnotation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)

	idStr := d.Id()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid id: %#v", idStr)
	}

	_, err = client.DeleteAnnotation(id)

	return err
}

func makeAnnotation(d *schema.ResourceData) (*gapi.Annotation, error) {
	idStr := d.Id()
	var id int64
	var err error
	if idStr != "" {
		id, err = strconv.ParseInt(idStr, 10, 64)
	}

	return &gapi.Annotation{
		ID:   id,
		Text: d.Get("text").(string),
	}, err
}
