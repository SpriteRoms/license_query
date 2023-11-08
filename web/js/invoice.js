function getParameterByName(name, url = window.location.href) {
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

(async function () {
    const go = new Go();
    let mod, inst;
    await WebAssembly.instantiateStreaming(fetch("js/invoice.wasm"), go.importObject).then((result) => {
        mod = result.module;
        inst = result.instance;
    }).catch((err) => {
        console.error(err);
    });
    await go.run(inst);
})()

async function INITED_EVENT() {
    const headingTemplate = document.getElementById('HeadingTemplate');
    const itemTemplate = document.getElementById('ItemTemplate');
    const licenseTable = document.getElementById('ItemTable');

    const activeCode = getParameterByName("activeCode")
    const license = queryLicense(activeCode)
    
    var licensePacksUrl = ``

    var pengingResourceNodes = []
    var pendingLicenseNodes = []

    await fetch("resources.json")
        .then(response => response.json())
        .then(json => {
            json.Tools.forEach(resource => {
                console.log(resource)
                var resourceGroupNode = headingTemplate.content.cloneNode(true);
                resourceGroupNode.getElementById("Left").textContent = resource.GroupName
                pengingResourceNodes.push(resourceGroupNode);

                resource.Resources.forEach(item => {
                    var resourceNode = itemTemplate.content.cloneNode(true);
                    resourceNode.getElementById("Left").textContent = item.Name
                    resourceNode.getElementById("Right").innerHTML = `<a href="${item.Link}">下载</a>`
                    pengingResourceNodes.push(resourceNode);
                })

            })
            licensePacksUrl = json.LicensePack
            return json
        })

    {
        var licenseNode = headingTemplate.content.cloneNode(true);
        licenseNode.getElementById("Left").textContent = `设备码: ${license.DeviceId}`
        licenseNode.getElementById("Right").textContent = `激活码: ${activeCode.substring(0, 5) + "..." + activeCode.substring(activeCode.length - 5)}`

        pendingLicenseNodes.push(licenseNode);

        license.Products.forEach(product => {
            var productNode = itemTemplate.content.cloneNode(true);
            productNode.getElementById("Left").textContent = `版本: ${product}`
            productNode.getElementById("Right").innerHTML = `<a href="${licensePacksUrl}">授权包ID: ${license.LicenseId.padStart(5, '0')}</a>`
            pendingLicenseNodes.push(productNode);
        })
    }

    pendingLicenseNodes.forEach(node => {
        licenseTable.appendChild(node);
    })
    pengingResourceNodes.forEach(node => {
        licenseTable.appendChild(node);
    })
}