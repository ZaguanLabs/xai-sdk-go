package files

import (
	"bytes"
	"context"
	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"net"
	"testing"
	"time"
)

type nativeFileServer struct {
	xaiv1.UnimplementedFilesServer
	t       *testing.T
	content []byte
}

func (s *nativeFileServer) UploadFile(stream grpc.ClientStreamingServer[xaiv1.UploadFileChunk, xaiv1.File]) error {
	first, err := stream.Recv()
	if err != nil {
		return err
	}
	if first.GetInit().GetName() != "test.bin" || first.GetInit().GetExpiresAfter() != 60 {
		s.t.Error(first)
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s.content = append(s.content, chunk.GetData()...)
	}
	return stream.SendAndClose(&xaiv1.File{Id: "file", Filename: "test.bin", Size: int64(len(s.content))})
}
func (s *nativeFileServer) RetrieveFile(_ context.Context, r *xaiv1.RetrieveFileRequest) (*xaiv1.File, error) {
	return &xaiv1.File{Id: r.FileId, Size: 123}, nil
}
func (s *nativeFileServer) RetrieveFileContent(r *xaiv1.RetrieveFileContentRequest, stream grpc.ServerStreamingServer[xaiv1.FileContentChunk]) error {
	if r.GetFormat() != xaiv1.DownloadFormat_DOWNLOAD_FORMAT_TEXT {
		s.t.Error(r)
	}
	for _, b := range [][]byte{[]byte("one"), {}, []byte("two")} {
		if err := stream.Send(&xaiv1.FileContentChunk{Data: b}); err != nil {
			return err
		}
	}
	return nil
}
func (s *nativeFileServer) ListFiles(_ context.Context, r *xaiv1.ListFilesRequest) (*xaiv1.ListFilesResponse, error) {
	if r.PaginationToken == nil || r.GetPaginationToken() != "next" || r.SortBy == nil || r.GetOrder() != xaiv1.FilesOrdering_DESCENDING {
		s.t.Error(r)
	}
	next := "last"
	return &xaiv1.ListFilesResponse{Data: []*xaiv1.File{{Id: "file"}}, PaginationToken: &next}, nil
}
func TestNativeFilesProtocol(t *testing.T) {
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	impl := &nativeFileServer{t: t}
	xaiv1.RegisterFilesServer(server, impl)
	go server.Serve(listener)
	defer server.Stop()
	conn, err := grpc.NewClient("passthrough:///files-test", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := NewGRPCClient(xaiv1.NewFilesClient(conn), nil)
	content := bytes.Repeat([]byte{0, 1, 2, 255}, (4<<20)/4)
	expiry := time.Minute
	file, err := client.Upload(context.Background(), bytes.NewReader(content), UploadOptions{Name: "test.bin", ExpiresAfter: &expiry})
	if err != nil {
		t.Fatal(err)
	}
	if file.Size != int64(len(content)) || !bytes.Equal(content, impl.content) {
		t.Fatal("upload corrupted")
	}
	file, err = client.Get(context.Background(), "file")
	if err != nil || file.Size != 123 {
		t.Fatal(file, err)
	}
	format := xaiv1.DownloadFormat_DOWNLOAD_FORMAT_TEXT
	reader, err := client.DownloadWithFormat(context.Background(), "file", &format)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil || string(data) != "onetwo" {
		t.Fatal(string(data), err)
	}
	result, err := client.List(context.Background(), &ListOptions{PaginationToken: "next", Order: xaiv1.FilesOrdering_DESCENDING})
	if err != nil || result.PaginationToken != "last" || len(result.Files) != 1 {
		t.Fatal(result, err)
	}
}
