# OTGSITEBUILDER

This is a static blog generator which it takes another different approach
for metadata. Instead of store metadata on posts (frontmatter), store them in a database.
So you use commands and/or paramaters to manipulate metadata related with
you content and you can use or write content with any tool or program.

The application has two modes: *builder* and *manager*; and always operate
with a folder with content.

With *manager* mode you operate with metadata and content options.
With *builder* mode you contruct the site with the content, metadata and content options.

## Commands

### Base

otgsitebuilder -mode (builder|manager) -content <path to folder>

### Builder

The command below is enough to build your website to the folder *output*, which can be configured.

otgsitebuilder -mode builder -content <path to folder>

### Manager
