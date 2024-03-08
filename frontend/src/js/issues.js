import "../css/home.css";
import "../css/issues.css";
import "../css/color-palette.css";
import data from "../database.json"; //get data from database



function openIssuesPage() {
    document.getElementById("page-contents").innerHTML = `
    <table id="issues-table">
        <thead>
            <tr>
                <th id="issue-column"><span class="table-header">Name</span><span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span></th>
                <th id="type-column"><span class="table-header">Type</span><span class="material-symbols-outlined" id="sort-on-type">swap_vert</span></th>
                <th id="risk-column"><span class="table-header">Risk level</span><span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span></th>
                <th id="risk-column"><span class="table-header">Information</span></th>
            </tr>
        </thead>
        <tbody>
        </tbody>
    </table>
    <script src="./src/js/issues-table.js" type="module" defer></script>
    `;

    // Fill the table with issues
    for (let i = 0; i < data.length; i++) {

        const body = document.querySelector('tbody')
        var row = `<tr>
        <td>${data[i].Name}</td>
        <td>${data[i].Type}</td>
        <td>${data[i].Risk}</td>
        <td>${data[i].Information}</td>`
        body.innerHTML += row;
    }


    // Sort the table 
    function sortTable(column) {
        console.log("you clicked on column header " + column);
    
        var table = document.getElementById("issues-table");
        var direction = 0;
    
        var swap = true;
        var count = 0;
        var rows;
        var swapTheseRows;
        while (swap) {
            swap = false;
            rows = table.rows;
            for (var i = 1; i < (rows.length-1); i++) {
                swapTheseRows = false;
                let elem1 = rows[i].getElementsByTagName("td")[column];
                let elem2 = rows[i+1].getElementsByTagName("td")[column];
            
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
}

document.getElementById("issues-button").addEventListener("click", () => openIssuesPage());

