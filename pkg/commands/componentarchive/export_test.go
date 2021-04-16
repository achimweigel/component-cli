// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package componentarchive_test

import (
	"context"
	"path/filepath"

	"github.com/gardener/component-spec/bindings-go/ctf"
	"github.com/mandelsoft/vfs/pkg/layerfs"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/projectionfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gardener/component-cli/pkg/commands/componentarchive"
	"github.com/gardener/component-cli/pkg/utils"
)

var _ = Describe("Export", func() {

	var testdataFs vfs.FileSystem

	BeforeEach(func() {
		baseFs, err := projectionfs.New(osfs.New(), "./testdata")
		Expect(err).ToNot(HaveOccurred())
		testdataFs = layerfs.New(memoryfs.New(), baseFs)
	})

	Context("From Filesystem", func() {

		It("should export a component archive from filesystem as tar file", func() {
			opts := &componentarchive.ExportOptions{
				ComponentArchivePath: "00-ca",
				OutputPath:           "ca.tar",
				OutputFormat:         ctf.ArchiveFormatTar,
			}

			Expect(opts.Run(context.TODO(), testdataFs)).To(Succeed())
			mediatype, err := utils.GetFileType(testdataFs, "ca.tar")
			Expect(err).ToNot(HaveOccurred())
			Expect(mediatype).ToNot(ContainSubstring("gzip"))
		})

		It("should export a component archive from filesystem as tar file", func() {
			opts := &componentarchive.ExportOptions{
				ComponentArchivePath: "00-ca",
				OutputPath:           "ca.tar.gz",
				OutputFormat:         ctf.ArchiveFormatTarGzip,
			}

			Expect(opts.Run(context.TODO(), testdataFs)).To(Succeed())
			mediatype, err := utils.GetFileType(testdataFs, "ca.tar.gz")
			Expect(err).ToNot(HaveOccurred())
			Expect(mediatype).To(Equal("application/x-gzip"))
		})

	})

	Context("From tar", func() {

		It("should export a component archive as tar file to filesystem", func() {
			opts := &componentarchive.ExportOptions{
				ComponentArchivePath: "00-ca",
				OutputPath:           "ca.tar",
				OutputFormat:         ctf.ArchiveFormatTar,
			}

			Expect(opts.Run(context.TODO(), testdataFs)).To(Succeed())

			opts = &componentarchive.ExportOptions{
				ComponentArchivePath: "ca.tar",
				OutputPath:           "ca",
				OutputFormat:         ctf.ArchiveFormatFilesystem,
			}

			Expect(opts.Run(context.TODO(), testdataFs)).To(Succeed())
			outputfileinfo, err := testdataFs.Stat("ca")
			Expect(err).ToNot(HaveOccurred())
			Expect(outputfileinfo.IsDir()).To(BeTrue(), "output filepath should be a directory")

			outputfileinfo, err = testdataFs.Stat(filepath.Join("ca", ctf.ComponentDescriptorFileName))
			Expect(err).ToNot(HaveOccurred())
			Expect(outputfileinfo.IsDir()).To(BeFalse(), "output filepath should be a directory")
		})

	})

})
