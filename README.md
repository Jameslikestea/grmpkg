## GRMPKG

GRMPKG is an immutable package hosting server for Go, it provides certainty in the eco-system where packages can disappear.

### Open Source

GRMPKG is licensed under the MIT License, you are free to deploy, modify and distribute the software in any way you like, however it is provided with no warranty.

### Design

GRMPKG conforms to the vgo standard, usually this is used for GOPROXY, however it means that it can be used to provide better, immutable package management in the cloud and in the datacenter. The vgo standard provides several endpoints and allows the download and install of go modules through the use of zip, go.mod and some info.

| Endpoint            | Purpose                                                                 |
| ------------------- | ----------------------------------------------------------------------- |
| .../@v/list         | List all of the versions available                                      |
| .../@v/version.info | Get information about v1, v1 can be replaced with any available version |
| .../@v/version.mod  | Download the go.mod file for the version                                |
| .../@v/version.zip  | Download the source code for the go module                              |
| .../@latest         | The latest version available, this will be done by **timestamp**        |

For now, if you require private modules, you should host this in your **own** VPC or Network

### Feature List

- [ ] Implementation of vgo implementation
- [ ] REST API to upload bundle
- [ ] Client to do upload
- [ ] Filesystem backend
- [ ] S3 backend
- [ ] SQL Service Information
