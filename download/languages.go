package download

import "context"

// Languages parses ISO-639 language codes, as used for alternate names in the file alternateNames.zip.
func (c *Client) Languages(ctx context.Context, callback func(parsed Language) error) error {
	header := true

	return c.downloadAndParseFile(ctx, "iso_languagecodes.txt", func(row []string) error {
		if header {
			header = false

			return nil
		}

		var parsed Language

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
