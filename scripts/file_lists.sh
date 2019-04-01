#!/bin/bash

# Top level folders whose name starts with a "number-".
dirs=$(find -E . -depth 1 -type d -regex '.*\/[0-9]+-.*')

last=""
lastDeps=""
for dir in $dirs; do

  # Nothing to do for the first pass.
  if [[ "${last}" == "" ]]; then
    last="${dir}"
    continue
  fi

  readme="${dir}/README.md"

  # Empty out the "File Changes" section from the README.md for this topic. It should
  # be the last part of the file. Start a new code fence as well.
  perl -0777 -pi -e "s/(## File Changes:).*/\1\n\n\`\`\`\n/s" "${readme}"

  # Find the changes between these two folders, format them with sed, and append them to the file.
  # Using git diff --raw prints output like this: (see man git-diff)
  # :100644 100644 0000000 0000000 M	./01-startup/main.go
  git diff --no-index --raw "${last}" "${dir}" | \
    sed -e "
    # Ignore changes to docs, examples, and modules.
    /README.md/d
    /doc.go/d
    /_examples/d
    /go.mod/d
    /go.sum/d

    # Remove the first 31 chars which hold the mode/sha1 stuff.
    s/^.\{31\}//

    # Make it easier to write the following sed scripts by converting tabs to
    # spaces. If your editor doesn't show it, there is a literal tab character
    # here between the first set of // characters
    s/	/ /g

    # Map the leading identifier to a full word. Renames are like R100 where
    # the files are 100% the same or R86 if the files are 86% similar. Also
    # match any trailing whitespace so we can control the output.
    s/^M/Modified/
    s/^A/Added   /
    s/^D/Deleted /
    s/^R100/Moved   /
    s/^R[[:digit:]]*/Moved+  /

    # Remove the topic directory names.
    s@${last}\/@@
    s@${dir}\/@@

    # For the case of moved files there will be two file names listed: the old
    # followed by the new. This monstrosity of a sed script matches up to the
    # second set of whitespace and replaces it with ' -> '. It uses matching
    # groups to preserve the first part of the line. It matches:
    # (1+ nonwhitespace 1+ space 1+ nonwhitespace) 1+ space (1+ nonwhitespace)
    s/^\([^ ]\{1,\} \{1,\}[^ ]\{1,\}\) \{1,\}\([^ ]\{1,\}\)/\1 -> \2/
    " >> "${readme}"

  # Close the code fence.
  echo -e -n "\`\`\`" >> "${readme}"

  # Generate another section about new dependencies by diffing the go.mod files
  # of each folder. Only include this section if there are differences.
  newDeps=$(diff "${last}/go.mod" "${dir}/go.mod" |\
    sed -e "
      # Filter to only dependency lines and only direct deps.
      /^[^><]/d
      /[()]/d
      /indirect/d
      /^[><][[:space:]]*\$/d

      # Use +- instead of >< for differences.
      s/^</-/
      s/^>/+/")
  if [[ "${newDeps}" != "" ]]; then
    echo -e -n "\n\n## Dependency Changes:\n\n\`\`\`\n${newDeps}\n\`\`\`" >> "${readme}"
  fi

  # Set the $last variable for the next pass of the loop.
  last=$dir
done;
