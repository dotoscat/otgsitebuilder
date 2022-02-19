// Copyright 2021 Oscar Triano Garcí­a

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
    "testing"
)

func TestContent(t *testing.T) {
    // const CORE_2 = 8
    content := OpenContent("testdata/content")
    t.Log("content", content)
    batchCh := content.GetPostsByCategory(ALL, 3)
    i := 0
    for batch := range batchCh {
        t.Log(i, batch)
        //go func(){
           for post := range batch.Posts() {
               t.Log("post:", post)
        }
        //}()
        i++
    }
    //<-done
}

