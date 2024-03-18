// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package datapipeline

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/datapipeline"
	"github.com/aws/aws-sdk-go/service/datapipeline/datapipelineiface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// []*SERVICE.Tag handling

// Tags returns datapipeline service tags.
func Tags(tags tftags.KeyValueTags) []*datapipeline.Tag {
	result := make([]*datapipeline.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &datapipeline.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from datapipeline service tags.
func KeyValueTags(ctx context.Context, tags []*datapipeline.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns datapipeline service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*datapipeline.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets datapipeline service tags in Context.
func setTagsOut(ctx context.Context, tags []*datapipeline.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates datapipeline service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn datapipelineiface.DataPipelineAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.DataPipeline)
	if len(removedTags) > 0 {
		input := &datapipeline.RemoveTagsInput{
			PipelineId: aws.String(identifier),
			TagKeys:    aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.RemoveTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.DataPipeline)
	if len(updatedTags) > 0 {
		input := &datapipeline.AddTagsInput{
			PipelineId: aws.String(identifier),
			Tags:       Tags(updatedTags),
		}

		_, err := conn.AddTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates datapipeline service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).DataPipelineConn(ctx), identifier, oldTags, newTags)
}
