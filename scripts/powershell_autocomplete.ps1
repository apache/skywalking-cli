# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

Register-ArgumentCompleter -Native -CommandName swctl -ScriptBlock {
      param($commandName, $wordToComplete, $cursorPosition)
      
      $match = $($(complete($wordToComplete)) -split " ")
      # Output matched commands one by one.
      for($i=0; $i -lt ($match.Length-1); $i+=1){
            Write-Output $match[$i]
      }
}
# Find all matching commands.
function complete($commands){
      # Get command line parameters.
      $parameters = $($commands -split " ")
      # Uncompleted parameters.
      $uncomplete = $parameters[-1]

      # Get the parameters before $uncomplete.
      $len = $parameters.Length-2
      $beforeCommands = $parameters[0..($len)]

      # Find the command prefixed with $uncomplete.
      $match = ""
      Invoke-Expression "$beforeCommands --auto_complete" | ForEach-Object {
            $flag = 1
            for ($i=0; $i -lt $uncomplete.Length; $i = $i +1){
                  if ($_[$i] -ne $uncomplete[$i]) { $flag = 0 }
            }
            if ($flag -eq 1) {  $match+="$_ " }
      }
      return $match
}

