package collections

import (
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

// fromProtoCollection converts a proto CollectionMetadata to a Collection.
func fromProtoCollection(pc *xaiv1.CollectionMetadata) *Collection {
	if pc == nil {
		return nil
	}

	c := &Collection{
		ID:             pc.CollectionId,
		Name:           pc.CollectionName,
		DocumentsCount: pc.DocumentsCount,
	}

	if pc.CreatedAt != nil {
		c.CreatedAt = pc.CreatedAt.AsTime()
	}

	return c
}

// fromProtoDocument converts a proto DocumentMetadata to a Document.
func fromProtoDocument(pd *xaiv1.DocumentMetadata) *Document {
	if pd == nil {
		return nil
	}

	d := &Document{
		Status:   pd.Status,
		ErrorMsg: pd.ErrorMessage,
	}

	// Extract file metadata
	if pd.FileMetadata != nil {
		fm := pd.FileMetadata
		d.FileID = fm.FileId
		d.Name = fm.Name
		d.SizeBytes = fm.SizeBytes
		d.ContentType = fm.ContentType
		d.Hash = fm.Hash

		if fm.CreatedAt != nil {
			d.CreatedAt = fm.CreatedAt.AsTime()
		}

		if fm.ExpiresAt != nil {
			d.ExpiresAt = fm.ExpiresAt.AsTime()
		}
	}

	// Convert fields
	if len(pd.Fields) > 0 {
		d.Fields = make(map[string]string, len(pd.Fields))
		for _, f := range pd.Fields {
			d.Fields[f.Key] = f.Value
		}
	}

	return d
}
