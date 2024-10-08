package main

import (
	"context"
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type HelloRoot struct {
	fs.Inode
}

func (r *HelloRoot) OnAdd(ctx context.Context) {
	ch := r.NewPersistentInode(ctx, &fs.MemRegularFile{
		Data: []byte("file.txt"),
		Attr: fuse.Attr{Mode: 0644},
	},
		fs.StableAttr{Ino: 2},
	)
	r.AddChild("file.txt", ch, false)
}

func (r *HelloRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}

var _ fs.NodeGetattrer = (*HelloRoot)(nil)
var _ fs.NodeOnAdder = (*HelloRoot)(nil)

func main() {
	listPod()
	// debug := flag.Bool("debug", false, "print debug data")
	// flag.Parse()
	// if len(flag.Args()) < 1 {
	// 	log.Fatal("Usage:\n hello [MOUNTPOINT]")
	// }
	// opts := &fs.Options{}
	// opts.Debug = *debug
	// server, err := fs.Mount(flag.Arg(0), &HelloRoot{}, opts)
	// if err != nil {
	// 	log.Fatalf("Mount failed: %v\n", err)
	// }
	// server.Wait()
}

func listPod() {
	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	clientset, _ := kubernetes.NewForConfig(config)
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Printf("%s  %s\n", pod.GetNamespace(), pod.GetName())
	}
}
