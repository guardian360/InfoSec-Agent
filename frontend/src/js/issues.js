import "../css/home.css";
import "../css/issues.css";
import "../css/color-palette.css";
import data from "../database.json"; //get data from database
import { openIssuePage } from "./issue.js";

export function openIssuesPage() {
    const pageContents = document.getElementById("page-contents");
    pageContents.innerHTML = `
    <table id="issues-table">
        <thead>
            <tr>
            <th id="issue-column"><span class="table-header">Name</span><span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span></th>
            <th id="type-column"><span class="table-header">Type</span><span class="material-symbols-outlined" id="sort-on-type">swap_vert</span></th>
            <th id="risk-column"><span class="table-header">Risk level</span><span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span></th>
            </tr>
        </thead>
        <tbody>
        </tbody>
    </table>
    `;

    var issues = []; // retrieve issues from tray application
    // dummy info
    issues = [
        { 
            "Id": "Windows defender", 
            "Result": ["Windows defender is disabled"],
            "ErrorMSG": null
        },
        { 
            "Id": "Camera and microphone access", 
            "Result": ["Something has access to camera", "Something has access to microphone"],
            "ErrorMSG": null
        }
    ];

    const tbody = pageContents.querySelector('tbody');

    issues.forEach(issue => {
        const currentIssue = data.find(element => element.Name === issue.Id);
        if (currentIssue) {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td class="issue-link">${currentIssue.Name}</td>
                <td>${currentIssue.Type}</td>
                <td>${currentIssue.Risk}</td>
            `;
            tbody.appendChild(row);
        }
    });

    const issueLinks = document.querySelectorAll(".issue-link");
    issueLinks.forEach((link, index) => {
        link.addEventListener("click", () => openIssuePage(issues[index].Id));
    });

    document.getElementById("sort-on-issue").addEventListener("click", () => sortTable(0));
    document.getElementById("sort-on-type").addEventListener("click", () => sortTable(1));
    document.getElementById("sort-on-risk").addEventListener("click", () => sortTable(2));
}

// Sort the table
function sortTable(column) {
    console.log("you clicked on column header " + column);

    const table = document.getElementById("issues-table");
    let direction = table.getAttribute("data-sort-direction");
    direction = direction === "ascending" ? "descending" : "ascending";

    const rows = Array.from(table.rows).slice(1); // Exclude header row
    rows.sort((a, b) => {
        const textA = a.cells[column].textContent.toLowerCase();
        const textB = b.cells[column].textContent.toLowerCase();
        if (direction === "ascending") {
            return textA.localeCompare(textB);
        } else {
            return textB.localeCompare(textA);
        }
    });

    while (table.rows.length > 1) {
        table.deleteRow(1); // Delete all rows except header
    }

    rows.forEach(row => {
        table.appendChild(row);
    });

    table.setAttribute("data-sort-direction", direction);
}

document.getElementById("issues-button").addEventListener("click", () => openIssuesPage());
