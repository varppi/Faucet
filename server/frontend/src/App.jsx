import {useState, useEffect} from 'react';

import {GetListener, GetConnected, GetFiles, PressedItem, DownloadItem, IsDownloading, KillAgent} from '../wailsjs/go/main/App'



function convBytes(x){
  const units = ['bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let l = 0, n = parseInt(x, 10) || 0;
  while(n >= 1024 && ++l){
      n = n/1024;
  }
  return(n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
}

function App() {
    var [waiting, setWaiting] = useState(<></>)
    let [connected, setConnected] = useState(false)
    var [preview, setPreview] = useState(false)
    var [downloadInProgress, setDownloadInProgress] = useState("")
    var [files, setFiles] = useState()
    var [previewUp, setPreviewUp] = useState(false)
    var [curDir, setCurDir] = useState("")
    
    useEffect(() => {
        var listener
        var connected2
        var changeBack = true
        var needRefresh = true
        
        GetListener().then((listener_) => {
            listener = listener_
            setWaiting(<h2>Waiting for connection on: {listener_}</h2>)
        }).then(() => {
            setInterval(() => {
                IsDownloading().then((downloading) => {
                    if (downloading != "") {
                        setDownloadInProgress(convBytes(downloading))
                    } else {
                        setDownloadInProgress("")
                    }
                })
                
                GetConnected().then((connected_) => {
                    if (connected2 && !connected_) {
                        setWaiting(<h2 style={{"color": "red"}}>Client disconnected!</h2>)
                        changeBack = false
                        setTimeout(() => {
                            changeBack = true
                        }, 4000);
                    } else if (!connected2 && (connected2 == connected_) && changeBack) {
                        setWaiting(<h2>Waiting for connection on: {listener}</h2>)
                    }
                    
                    connected2 = connected_
                    setConnected(connected_)
                    
                    
                    function newItem(name, val) {
                        function handleClick() {
                            needRefresh = true
                            PressedItem(name).then((preview) => {
                                setPreview(preview)
                            })
                        }
                        function builder(name, icon, val) {
                            if (val == "d") {
                                if (name == "../") {
                                    return  <li><button onClick={()=>{handleClick()}} style={{"color": "dodgerblue", "letterSpacing": "2px", "textIndent": "5px", "fontWeight": "700"}}>{name}</button></li>
                                }
                                return <li><button onClick={()=>{handleClick()}}>{icon}{name}</button><button onClick={()=>{DownloadItem(name)}} className='download'>Download</button></li>
                            } else {
                                val = convBytes(val)
                                return <li><button onClick={()=>{handleClick()}}>{icon}{name}</button><span className="size"> {val}</span><button onClick={()=>{DownloadItem(name)}} className='download'>Download</button></li>
                            }
                        }
                        if (val == "d") {
                            if (name == "../") {
                                return builder(name, "", val)
                            }
                            return builder(name, "📁", val)
                        } else {
                            if (name.endsWith(".exe")) {
                                return builder(name, "⚙️", val)
                            }
                            if (name.endsWith(".conf")) {
                                return builder(name, "🔧", val)
                            }
                            if (name.endsWith(".lnk")) {
                                return builder(name, "🔗", val)
                            }
                            if (name.endsWith(".sql") || name.endsWith(".db") ||name.endsWith(".sqlite3") ||name.endsWith(".sqlite") ||name.endsWith(".mdb")) {
                                return builder(name, "💎", val)
                            }
                            if (name.endsWith(".html") || name.endsWith(".php") || name.endsWith(".jsx") || name.endsWith(".js") || name.endsWith(".css") || name.endsWith(".htm")) {
                                return builder(name, "🌐", val)
                            }
                            if (name.endsWith(".txt")) {
                                return builder(name, "📄", val)
                            }
                            if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png") || name.endsWith(".webp")) {
                                return builder(name, "🖼", val)
                            }
                            return builder(name, "❔", val)
                        }
                        
                    }
                    
                    if (connected_ && needRefresh) {
                        GetFiles().then((files) => {
                            var allHtml = []
                            allHtml.push(newItem("../", "d"))
                            
                            var dirsHtml = []
                            var filesHtml = []

                            for (let [key,val] of Object.entries(files)) {
                                if (key == "/curdir/") {
                                    setCurDir(val)
                                } else if (val == "d") {
                                    dirsHtml.push(newItem(key, "d"))
                                } else {
                                    filesHtml.push(newItem(key, val))
                                }
                            }
                            filesHtml.sort()
                            dirsHtml.sort()
                            filesHtml.forEach((val) => {
                                allHtml.push(val)
                            })
                            dirsHtml.forEach((val) => {
                                allHtml.push(val)
                            })
                            needRefresh = false
                            setFiles(allHtml)
                        })
                    }
                })
            }, 500)
        })
    }, [])
    
    var previewBox = <div className="previewBox" style={{"maxWidth": (previewUp ? "100%" : "500px")}}>
        <div className="buttons">
            <button className="btn previewbutton" onClick={()=>{setPreviewUp(!previewUp)}}>^</button>
            <button className="btn previewbutton" onClick={()=>{setPreview(null)}}>X</button>
        </div>
        <pre>
            {preview}
        </pre>
    </div>

    var filesBox = <div className="bigmargin d-flex flex">
        {(preview && !previewUp) ? previewBox: null}
        <div className="box">
            <p style={{"fontSize": "12px"}}>{curDir}</p>
            {files}
        </div>
    </div>
    
    var connectedHt = <div className="d-flex justify-content-between mx-2">
            <h2 className="maincolor" id="connected">Connected!</h2>
            <button onClick={()=>{KillAgent()}} id="killagent" className="btn">Kill Agent</button>
        </div>

    var downloading = <h5 id="downloading" className={downloadInProgress ? "": "invisible"}>download progress: {downloadInProgress}</h5>
    return (
        <>
        {connected ? connectedHt: waiting}
        {downloading}
        {(preview && previewUp) ? previewBox: null}
        {connected ? filesBox: null}
        </>
    )
}

export default App
