package main

import "strings"

func ossLicenses() []OSSLicense {
	return []OSSLicense{
		{Name: "ClipForVRChat", License: "MIT", Copyright: "Copyright (c) 2026 hatolife", URL: "https://github.com/hatolife/ClipForVRChat", Text: mitLicenseText("Copyright (c) 2026 hatolife")},
		{Name: "@vitejs/plugin-vue", License: "MIT", Copyright: "Copyright (c) 2019-present, Yuxi (Evan) You and Vite contributors", URL: "https://github.com/vitejs/vite-plugin-vue", Text: mitLicenseText("Copyright (c) 2019-present, Yuxi (Evan) You and Vite contributors")},
		{Name: "Vite", License: "MIT", Copyright: "Copyright (c) 2019-present, VoidZero Inc. and Vite contributors", URL: "https://github.com/vitejs/vite", Text: mitLicenseText("Copyright (c) 2019-present, VoidZero Inc. and Vite contributors")},
		{Name: "Vue.js", License: "MIT", Copyright: "Copyright (c) 2018-present, Yuxi (Evan) You", URL: "https://github.com/vuejs/core", Text: mitLicenseText("Copyright (c) 2018-present, Yuxi (Evan) You")},
		{Name: "Wails", License: "MIT", Copyright: "Copyright (c) 2018-Present Lea Anthony", URL: "https://github.com/wailsapp/wails", Text: mitLicenseText("Copyright (c) 2018-Present Lea Anthony")},
		{Name: "cloudflare/circl", License: "BSD-3-Clause", Copyright: "Copyright (c) 2019 Cloudflare; Copyright (c) 2009 The Go Authors", URL: "https://github.com/cloudflare/circl", Text: circlLicenseText},
		{Name: "flock", License: "BSD-2-Clause", Copyright: "Copyright (c) 2018-2025, The Gofrs; Copyright (c) 2015-2020, Tim Heckman", URL: "https://github.com/gofrs/flock", Text: bsd2ClauseLicenseText("Copyright (c) 2018-2025, The Gofrs\nCopyright (c) 2015-2020, Tim Heckman")},
		{Name: "go-ansi-parser", License: "MIT", Copyright: "Copyright (c) 2021-Present Lea Anthony", URL: "https://github.com/leaanthony/go-ansi-parser", Text: mitLicenseText("Copyright (c) 2021-Present Lea Anthony")},
		{Name: "go-arg", License: "BSD-2-Clause", Copyright: "Copyright (c) 2015, Alex Flint", URL: "https://github.com/alexflint/go-arg", Text: bsd2ClauseLicenseText("Copyright (c) 2015, Alex Flint")},
		{Name: "go-scalar", License: "BSD-2-Clause", Copyright: "Copyright (c) 2015, Alex Flint", URL: "https://github.com/alexflint/go-scalar", Text: bsd2ClauseLicenseText("Copyright (c) 2015, Alex Flint")},
		{Name: "go-webview2", License: "MIT", Copyright: "Copyright (c) 2020 John Chadwick; Some portions Copyright (c) 2017 Serge Zaitsev", URL: "https://github.com/wailsapp/go-webview2", Text: mitLicenseText("Copyright (c) 2020 John Chadwick\nSome portions Copyright (c) 2017 Serge Zaitsev")},
		{Name: "golang.design/x/clipboard", License: "MIT", Copyright: "Copyright (c) 2021 Changkun Ou", URL: "https://github.com/golang-design/clipboard", Text: mitLicenseText("Copyright (c) 2021 Changkun Ou <contact@changkun.de>")},
		{Name: "golang.org/x/crypto", License: "BSD-3-Clause", Copyright: "Copyright 2009 The Go Authors", URL: "https://cs.opensource.google/go/x/crypto", Text: goBSD3ClauseLicenseText("Copyright 2009 The Go Authors.", "Google LLC")},
		{Name: "golang.org/x/image", License: "BSD-3-Clause", Copyright: "Copyright 2009 The Go Authors", URL: "https://cs.opensource.google/go/x/image", Text: goBSD3ClauseLicenseText("Copyright 2009 The Go Authors.", "Google LLC")},
		{Name: "golang.org/x/sys", License: "BSD-3-Clause", Copyright: "Copyright 2009 The Go Authors", URL: "https://cs.opensource.google/go/x/sys", Text: goBSD3ClauseLicenseText("Copyright 2009 The Go Authors.", "Google LLC")},
		{Name: "golang.org/x/text", License: "BSD-3-Clause", Copyright: "Copyright 2009 The Go Authors", URL: "https://cs.opensource.google/go/x/text", Text: goBSD3ClauseLicenseText("Copyright 2009 The Go Authors.", "Google LLC")},
		{Name: "golang.org/x/xerrors", License: "BSD-3-Clause", Copyright: "Copyright (c) 2019 The Go Authors", URL: "https://cs.opensource.google/go/x/xerrors", Text: goBSD3ClauseLicenseText("Copyright (c) 2019 The Go Authors. All rights reserved.", "Google Inc.")},
		{Name: "gozxing", License: "MIT", Copyright: "Copyright (c) 2018 Daisuke MAKIUCHI", URL: "https://github.com/makiuchi-d/gozxing", Text: mitLicenseText("Copyright (c) 2018 Daisuke MAKIUCHI (MakKi; makki_d)")},
		{Name: "imaging", License: "MIT", Copyright: "Copyright (c) 2012 Grigory Dryapak", URL: "https://github.com/disintegration/imaging", Text: mitLicenseText("Copyright (c) 2012 Grigory Dryapak")},
		{Name: "pkg/errors", License: "BSD-2-Clause", Copyright: "Copyright (c) 2015, Dave Cheney", URL: "https://github.com/pkg/errors", Text: bsd2ClauseLicenseText("Copyright (c) 2015, Dave Cheney <dave@cheney.net>")},
		{Name: "ProtonMail/go-crypto", License: "BSD-3-Clause", Copyright: "Copyright (c) 2009 The Go Authors", URL: "https://github.com/ProtonMail/go-crypto", Text: goBSD3ClauseLicenseText("Copyright (c) 2009 The Go Authors. All rights reserved.", "Google Inc.")},
		{Name: "rivo/uniseg", License: "MIT", Copyright: "Copyright (c) 2019 Oliver Kuederle", URL: "https://github.com/rivo/uniseg", Text: mitLicenseText("Copyright (c) 2019 Oliver Kuederle")},
		{Name: "slicer", License: "MIT", Copyright: "Copyright (c) 2019 Lea Anthony", URL: "https://github.com/leaanthony/slicer", Text: mitLicenseText("Copyright (c) 2019 Lea Anthony")},
		{Name: "Spout2", License: "BSD-2-Clause", Copyright: "Copyright (c) 2016-2025, Lynn Jarvis", URL: "https://github.com/leadedge/Spout2", Text: spout2LicenseText},
		{Name: "u", License: "MIT", Copyright: "Copyright (c) 2023-Present Lea Anthony", URL: "https://github.com/leaanthony/u", Text: mitLicenseText("Copyright (c) 2023-Present Lea Anthony")},
		{Name: "AvatarBeacon / YL-ATG", License: "MIT", Copyright: "Copyright (c) 2024 YozoraKurage", URL: "https://github.com/YozoraKurage/YL-ATG", Text: avatarBeaconYLATGLicenseText, Group: "avatarBeacon"},
	}
}

func mitLicenseText(copyright string) string {
	return strings.TrimSpace("MIT License\n\n" + strings.TrimSpace(copyright) + "\n\n" + mitLicenseTerms)
}

func bsd2ClauseLicenseText(copyright string) string {
	return strings.TrimSpace(strings.TrimSpace(copyright) + "\nAll rights reserved.\n\n" + bsd2ClauseLicenseTerms)
}

func goBSD3ClauseLicenseText(copyright string, neitherName string) string {
	return strings.TrimSpace(strings.TrimSpace(copyright) + "\n\n" + strings.ReplaceAll(goBSD3ClauseLicenseTerms, "{NEITHER_NAME}", neitherName))
}

const mitLicenseTerms = `Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`

const bsd2ClauseLicenseTerms = `Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

const goBSD3ClauseLicenseTerms = `Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of {NEITHER_NAME} nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

const circlLicenseText = `Copyright (c) 2019 Cloudflare. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Cloudflare nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

========================================================================

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

const spout2LicenseText = `Spout2 / SpoutLibrary

Copyright (c) 2016-2025, Lynn Jarvis. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`

const avatarBeaconYLATGLicenseText = `AvatarBeacon includes material derived from YL-ATG_ForAvatar V0.0.3 by YozoraKurage.

Upstream project:
https://github.com/YozoraKurage/YL-ATG

Modifications in AvatarBeacon:
- Renamed public asset path from Assets/YozoLab/YL-ATG_ForAvatar to Assets/PoppoWorks/AvatarBeacon.
- Renamed position parameters from ATG/p/* to avatar_beacon/coord/*.
- Renamed forward/rotation-vector parameters from ATG/r/* to avatar_beacon/forward/*.
- Removed auxiliary ATG/SaveObject and debug-only menu parameters that are not used for basis reconstruction.
- Kept the tracking Bone Proxy target on Head for both position and forward/yaw.
- Removed the separate HeadForwardAnchor transform and reused point for the forward/yaw sensor graph.
- Removed the visual-only arrow mesh/material assets.
- Normalized near-zero serialized Transform values in the prefab while preserving Contact/Constraint values used by the sensor graph.
- Adjusted prefab naming to AvatarBeacon_main and AvatarBeacon_12.

The MIT License (MIT)

Copyright (c) 2024 YozoraKurage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
