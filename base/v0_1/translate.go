// Copyright 2019 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.)

package v0_1

import (
	"net/url"

	"github.com/coreos/fcct/translate"

	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/coreos/vcontext/path"
	"github.com/coreos/vcontext/report"
	"github.com/vincent-petithory/dataurl"
)

// ToIgn3_0 translates the config to an Ignition config. It also returns the set of translations
// it did so paths in the resultant config can be tracked back to their source in the source config.
func (c Config) ToIgn3_0() (types.Config, translate.TranslationSet, report.Report) {
	ret := types.Config{}
	tr := translate.NewTranslator("yaml", "json")
	tr.AddCustomTranslator(translateIgnition)
	tr.AddCustomTranslator(translateFile)
	tr.AddCustomTranslator(translateDirectory)
	tr.AddCustomTranslator(translateLink)
	translations := tr.Translate(&c, &ret)
	return ret, translations, report.Report{}
}

func translateIgnition(from Ignition) (to types.Ignition, tm translate.TranslationSet) {
	tr := translate.NewTranslator("yaml", "json")
	to.Version = types.MaxVersion.String()
	tm = translate.Prefixed(tr, "config", &from.Config, &to.Config)
	translate.MergeP(tr, tm, "security", &from.Security, &to.Security)
	translate.MergeP(tr, tm, "timeouts", &from.Timeouts, &to.Timeouts)
	return
}

func translateFile(from File) (to types.File, tm translate.TranslationSet) {
	tr := translate.NewTranslator("yaml", "json")
	tr.AddCustomTranslator(translateFileContents)
	tm = translate.Prefixed(tr, "group", &from.Group, &to.Group)
	translate.MergeP(tr, tm, "user", &from.User, &to.User)
	translate.MergeP(tr, tm, "append", &from.Append, &to.Append)
	translate.MergeP(tr, tm, "contents", &from.Contents, &to.Contents)
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	to.Mode = from.Mode
	tm.AddIdentity("overwrite", "path", "mode")
	return
}

func translateFileContents(from FileContents) (to types.FileContents, tm translate.TranslationSet) {
	tr := translate.NewTranslator("yaml", "json")
	tm = translate.Prefixed(tr, "verification", &from.Verification, &to.Verification)
	to.Source = from.Source
	to.Compression = from.Compression
	tm.AddIdentity("source", "compression")
	if from.Inline != nil {
		src := (&url.URL{
			Scheme: "data",
			Opaque: "," + dataurl.EscapeString(*from.Inline),
		}).String()
		to.Source = &src
		tm.AddTranslation(path.New("yaml", "inline"), path.New("json", "source"))
	}
	return
}

func translateDirectory(from Directory) (to types.Directory, tm translate.TranslationSet) {
	tr := translate.NewTranslator("yaml", "json")
	tm = translate.Prefixed(tr, "group", &from.Group, &to.Group)
	translate.MergeP(tr, tm, "user", &from.User, &to.User)
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	to.Mode = from.Mode
	tm.AddIdentity("overwrite", "path", "mode")
	return
}

func translateLink(from Link) (to types.Link, tm translate.TranslationSet) {
	tr := translate.NewTranslator("yaml", "json")
	tm = translate.Prefixed(tr, "group", &from.Group, &to.Group)
	translate.MergeP(tr, tm, "user", &from.User, &to.User)
	to.Target = from.Target
	to.Hard = from.Hard
	to.Overwrite = from.Overwrite
	to.Path = from.Path
	tm.AddIdentity("target", "hard", "overwrite", "path")
	return
}
