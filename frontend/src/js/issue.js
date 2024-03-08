import data from "../database.json"; //get data from database
import "../css/issues.css";
import { openIssuesPage } from "./issues.js";

export function openIssuePage(issueId) {
    console.log("opened issue page: " + issueId);

    const currentIssue = data.find((element) => element.Name = issueId);
    document.getElementById("page-contents").innerHTML = `
        <h1 id="issue-name">${currentIssue.Name}</h1>
        <div id="issue-information">
            <h2>Information</h2>
            <p>${currentIssue.Information}</p>
            <h2>Solution</h2>
            <p>${currentIssue.Solution}</p>
        </div>
        <div id="back-button">Back to issues overview</div>
    `;

    document.getElementById("back-button").addEventListener("click", () => openIssuesPage());
}