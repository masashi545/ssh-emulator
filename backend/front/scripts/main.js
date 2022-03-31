let crypto;
try {
    crypto = require('crypto');
} catch (err) {
    console.log('crypto support is disabled!');
}

let buffer;
try {
    buffer = require('buffer');
} catch (err) {
    console.log('crypto support is disabled!');
}

const kexButton = document.getElementById("id-kex-button")
kexButton.addEventListener("click", (event) => {
    const userName = document.getElementById("id-username").value
    const kexAlgo = document.getElementById("id-kex-algorithm").value

    var dhGroup = ''
    if (kexAlgo == "diffie-hellman-group1-sha1") {
        dhGroup = 'modp2'
    } else { // kexAlgo == "diffie-hellman-group14-sha1"
        dhGroup = 'modp14'
    }


    //var alice = "3"
    //var aliceKey = "8"
    const alice = crypto.createDiffieHellmanGroup(dhGroup)
    const aliceKey = alice.generateKeys().toString('hex').toUpperCase()
    document.getElementById("id-alice-encrypted").value = aliceKey
    console.log("AliceLength: " + aliceKey.length * 4 + " bit")

    var shareKeyAlice = ''

    var bobKey = ''
    var shareKeyBob = ''

    // DH鍵交換を行うHTTP Request
    const req = new XMLHttpRequest() // XMLHttpRequest...ブラウザ上でサーバーとHTTP通信を行うためのAPI
    req.addEventListener("load", () => {

        const response = JSON.parse(req.responseText) // JSON形式 -> Object
        console.log(response)

        const res = Object.values(response)
        bobKey = res[0]
        shareKeyBob = res[3]
        console.log("BobLength: " + bobKey.length * 4 + " bit")

        document.getElementById("id-bob-encrypted").value = bobKey
        document.getElementById("id-shared-key-host").value = shareKeyBob

        shareKeyAlice = alice.computeSecret(Uint8Array.from(Buffer.from(bobKey.toLowerCase(), "hex")))
        document.getElementById("id-shared-key-user").value = shareKeyAlice.toString('hex').toUpperCase()


    })
    req.open(
        "POST", // HTTP Request Method 
        "http://localhost:8080/ssh-kex" // HTTP Request URL 
    )
    req.setRequestHeader("Context-Type", "application/json") // HTTP Request Header
    const httpBody = {
        User: userName,
        KexAlgo: kexAlgo,
        PubKey: aliceKey
    }
    req.send(JSON.stringify(httpBody))

}
)