Type: `xelon-template`

The Xelon template data source provides information about a template in Xelon HQ.

-> **Note:** Data sources is a feature exclusively available to HCL2 templates.

Basic example of usage:

```hcl
data "xelon-template" "basic-example" {
  name        = "ubuntu-24-base"
  most_recent = true
}

# usage example of the data source output
locals {
  template_id = data.xelon-template.basic-example.id
}
```

## Configuration Reference

Configuration options are organized below into two categories: required and optional. Within each category, the
available options are alphabetized and described.

**Required:**

<!-- Code generated from the comments of the TemplateConfig struct in internal/datasource/data_template.go; DO NOT EDIT MANUALLY -->

- `name` (string) - The template name.

<!-- End of code generated from the comments of the TemplateConfig struct in internal/datasource/data_template.go; -->


**Optional:**

<!-- Code generated from the comments of the TemplateConfig struct in internal/datasource/data_template.go; DO NOT EDIT MANUALLY -->

- `most_recent` (bool) - If true, the most recent OS template will be returned. If false (default),
  an error will be returned if more than one template matches the filters.

<!-- End of code generated from the comments of the TemplateConfig struct in internal/datasource/data_template.go; -->


## Output Data

<!-- Code generated from the comments of the TemplateDatasourceOutput struct in internal/datasource/data_template.go; DO NOT EDIT MANUALLY -->

- `id` (string) - The ID of the template.

- `name` (string) - The name of the template.

- `creation_date` (string) - The date of creation of the template.

<!-- End of code generated from the comments of the TemplateDatasourceOutput struct in internal/datasource/data_template.go; -->


## Authentication

This data source uses the same authentication method as the main Xelon plugin to connect to Xelon HQ. See the of
the [authentication section](/packer/integrations/Xelon-AG/xelon#authentication) plugin’s documentation for details.
