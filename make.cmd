call :"%~1"
exit /b

:""
:"all"
    go fmt || exit /b
    go build || exit /b
    pushd cmd\enexToHtml || exit /b
    go fmt && go build
    popd
    exit /b

:"clean"
    del *.png *.html
    exit /b

:"html"
    for %%I in (*.enex) do cmd\enexToHtml\enexToHtml.exe "%%~I"
    exit /b

:"md"
    for %%I in (*.enex) do cmd\enexToHtml\enexToHtml.exe -markdown "%%~I"
    exit /b

:"shrink"
    for %%I in (*.enex) do cmd\enexToHtml\enexToHtml.exe -shrink-markdown "%%~I"
    exit /b
