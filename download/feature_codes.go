package download

import (
	"context"
	"fmt"
)

// FeatureCodes parses name and description for feature classes and feature codes from a featureCodes_xx.txt file.
func (c *Client) FeatureCodes(ctx context.Context, language string) (Iterator[Feature], error) {
	res, err := c.downloadAndParseFile(ctx, fmt.Sprintf("featureCodes_%s.txt", language))

	return withUnmarshalRows[Feature](res), err
}
