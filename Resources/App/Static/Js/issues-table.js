var issues = []; // retrieve issues and fill issue table

// dummy info
issues = [
    { issue: "Windows defender disabled", type: "Security", risklevel: "high" },
    { issue: "Windows password", type: "Privacy", risklevel: "low" }
];


for (let i = 0; i < issues.length; i++) {
    var table = document.getElementById("issues-table");
    var row = table.insertRow(i+1);
    var cell1 = row.insertCell(0);
    var cell2 = row.insertCell(1);
    var cell3 = row.insertCell(2);
    cell1.innerHTML = `<a href="/home/">` + issues[i].issue + `</a>`;
    cell2.innerHTML = issues[i].type;
    cell3.innerHTML = issues[i].risklevel;
}