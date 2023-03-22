call :"%~1"
exit /b

:""
:"all"
    go fmt || exit /b
    go build || exit /b
    pushd cmd\enexToHtml || exit /b
    go fmt && go build -o ../..
    popd
    exit /b

:"clean"
    del *.png *.html
    exit /b

:"html"
    for %%I in (*.enex) do enexToHtml.exe "%%~I"
    exit /b

:"md"
    for %%I in (*.enex) do enexToHtml.exe -markdown "%%~I"
    exit /b

:"shrink"
    for %%I in (*.enex) do enexToHtml.exe -shrink-markdown "%%~I"
    exit /b

:"embed"
    for %%I in (*.enex) do enexToHtml.exe -embed "%%~I"
    exit /b
