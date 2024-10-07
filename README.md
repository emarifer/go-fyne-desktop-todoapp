<div align="center">

# fToDo App

<img src="assets/34.svg" width="45%">

<hr />

<p style="margin-bottom: 16px;">
    Golang desktop todo app using Fyne UI framework.
</p>

> 🚧 This is a work in progress and therefore you should expect that the
> application may not have all the features at this moment.

<br />

![GitHub License](https://img.shields.io/github/license/emarifer/go-fyne-desktop-todoapp) ![Static Badge](https://img.shields.io/badge/Go-%3E=1.23-blue)

</div>

<hr />

### Features 🚀

- [x] **Using the [Fyne](https://fyne.io/) UI framework:** `Fyne` is a framework that allows us to build native cross-platform graphical interfaces. Its documentation is well suited for beginners, the resulting binaries are small and allow you to package all your assets, and compile quickly (cross-compiling is possible). Although it is very popular and much easier to use than libraries like GTK-4 (binding for Go), with not too much code to do basic things, it may fall short when we want to build more complex interfaces. While Fyne is rapidly evolving and this situation may change 🙏…
- [x] **Using the [cloverDB](https://github.com/ostafen/clover) database:** It is a lightweight, document-oriented NoSQL database. Although we could have also used `SQLite`, this way we give a chance to this unusual database that, by not using SQL, is extraordinarily easy to use and as its creators say: "`CloverDB` has been written for being easily maintainable. As such, it trades performance with simplicity, and is not intended to be an alternative to more performant databases such as `MongoDB` or `MySQL`. However, there are projects where running a separate database server may result overkilled, and, for simple queries, network delay may be the major performance bottleneck. For such scenarios, `CloverDB` may be a more suitable alternative".

---