package fs


// self defined http error status
const (
    StatusOK    = iota
    StatusFileNotFound
    StatusMakeDirFail
    StatusDirectoryNotAFile
    StatusCopyFromLocalToLocal
    StatusDestShouldBeDirectory
    StatusOnlySupportFiles
    StatusBadFileSize
    StatusDirectoryAlreadyExist
    StatusBadChunkSize
    StatusShouldBePfsPath
    StatusNotEnoughArgs
    StatusInvalidArgs
    StatusUnAuthorized
    StatusJSONErr
    StatusCannotDelDirectory
    StatusAlreadyExist
    StatusBadPath
    OpenFileErr
    CreateFileErr
    SeekFileErr
    CopyFileErr
    ReadFileErr
    BadRawQueryURL
    EncodeJSONErr
    ParseStrToIntErr
    RemoveFileErr
)


var statusText = map[int]string {
    StatusOK : "Status OK",
    // StatusFileNotFound is a error string of that there is no file or directory.
    StatusFileNotFound : "no such file or directory",
    StatusMakeDirFail : "can't create directory",
    // StatusDirectoryNotAFile is a error string of that the destination should be a file.
    StatusDirectoryNotAFile : "should be a file not a directory",
    // StatusCopyFromLocalToLocal is a error string of that this system does't support copy local to local.
    StatusCopyFromLocalToLocal : "don't support copy local to local",
    // StatusDestShouldBeDirectory is a error string of that destination shoule be a directory.
    StatusDestShouldBeDirectory : "dest should be a directory",
    // StatusOnlySupportFiles is a error string of that the system only support upload or download files not directories.
    StatusOnlySupportFiles : "only support upload or download files not directories",
    // StatusBadFileSize is a error string of that the file size is no valid.
    StatusBadFileSize : "bad file size",
    // StatusDirectoryAlreadyExist is a error string of that the directory is already exist.
    StatusDirectoryAlreadyExist : "directory already exist",
    // StatusBadChunkSize is a error string of that the chunksize is error.
    StatusBadChunkSize : "chunksize error",
    // StatusShouldBePfsPath is a error string of that a path should be a pfs path.
    StatusShouldBePfsPath : "fs path should be begin with /",
    // StatusNotEnoughArgs is a error string of that there is not enough arguments.
    StatusNotEnoughArgs : "not enough arguments",
    // StatusInvalidArgs is a error string of that arguments are not valid.
    StatusInvalidArgs : "invalid arguments",
    // StatusUnAuthorized is a error string of that what you request should have authorization.
    StatusUnAuthorized : "what you request is unauthorized",
    // StatusJSONErr is a error string of that the system parses json error.
    StatusJSONErr : "parse json error",
    EncodeJSONErr : "encode json error",
    // StatusCannotDelDirectory is a error string of that what you input can't delete a directory.
    StatusCannotDelDirectory : "can't del directory",
    // StatusAlreadyExist is a error string of that the destination is already exist.
    StatusAlreadyExist : "already exist",
    // StatusBadPath is a error string of that the form of path is not correct.
    StatusBadPath : "the path is not correct",
    OpenFileErr: "failed to open file",
    CreateFileErr: "failed to create file",
    SeekFileErr: "failed to seek file",
    CopyFileErr: "copy file failed",
    ReadFileErr: "read file failed",
    BadRawQueryURL: "bad raw query url",
    ParseStrToIntErr: "parse string to int error",
    RemoveFileErr: "can't remove file",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
    return statusText[code]
}


