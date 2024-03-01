import "../css/home.css";
import "../css/issues.css";
import "../css/color-palette.css";

var issues = []; // retrieve issues and fill issue table

// dummy info
issues = [
    { issue: "Windows defender disabled", type: "Security", risklevel: 3 },
    { issue: "Windows password", type: "Privacy", risklevel: 1 }
];

// Fill the table with issues
for (let i = 0; i < issues.length; i++) {
    var table = document.getElementById("issues-table");
    var row = table.insertRow(i+1);
    var cell1 = row.insertCell(0);
    var cell2 = row.insertCell(1);
    var cell3 = row.insertCell(2);
    cell1.innerHTML = `<a href="/issue.html" class="issue-link">` + issues[i].issue + `</a>`;
    cell2.innerHTML = issues[i].type;
    if (issues[i].risklevel == 0) cell3.innerHTML = "safe";
    else if (issues[i].risklevel == 1) cell3.innerHTML = "low";
    else if (issues[i].risklevel == 2) cell3.innerHTML = "medium";
    else if (issues[i].risklevel == 3) cell3.innerHTML = "high";
    else cell3.innerHTML = "unknown";
}

// Sort the table 
function sortTable(column) {
    console.log("you clicked on column header " + column);

    var table = document.getElementById("issues-table");
    var direction = 0;

    var swap = true;
    var count = 0;
    while (swap) {
        swap = false;
        rows = table.rows;
        for (i = 1; i < (rows.length-1); i++) {
            swapTheseRows = false;
            elem1 = rows[i].getElementsByTagName("td")[column];
            elem2 = rows[i+1].getElementsByTagName("td")[column];

            if (direction == 0 && elem1.innerHTML.toLowerCase() > elem2.innerHTML.toLowerCase()) {
                swapTheseRows = true;
                break;
            }
            else if (direction == 1 && elem1.innerHTML.toLowerCase() < elem2.innerHTML.toLowerCase()) {
                swapTheseRows = true;
                break;
            }
        }
        if (swapTheseRows) {
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            swap = true;
            count++;
        }
        else if (direction == 0 && count == 0) {
            direction = 1;
            swap = true;
        }
    }
}

// When clicking on the symbols next to the column headers, the table is sorted according to that column
document.getElementById("sort-on-issue").addEventListener("click", () => sortTable(0));
document.getElementById("sort-on-type").addEventListener("click", () => sortTable(1));
document.getElementById("sort-on-risk").addEventListener("click", () => sortTable(2));