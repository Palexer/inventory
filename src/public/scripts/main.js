// set/update delete/edit functions
setDeleteEditFuncs()

// print button
document.getElementById("print-button").onclick = function () {window.print()}


// create frontend undocache
let undocache = new Array()

// undo button
document.getElementById("undo-button").onclick = function () {
	if (undocache.length < 1) {
		alert("No elements in undo cache")
		return
	}

	if (confirm("Do you want to restore the last removed item?")) {
		// undo request for backend
		let xhr = new XMLHttpRequest()
		xhr.open("POST", "/undo", true)
		xhr.send()

		// add row back on frontend
		let row = document.getElementById("main-table").insertRow(-1)
		for (i = 0; i < document.getElementsByTagName("th").length + 2; i++) {
			row.insertCell(i).innerHTML = undocache[undocache.length - 1].children[i].innerHTML
		}

		row.getElementsByTagName("td")[0].innerHTML = rows.length - 1

		// remove last element from cache
		undocache.pop()
		setDeleteEditFuncs()
	}
}


// sort table
/**
 * Sorts a HTML table.
 *
 * @param {HTMLTableElement} table The table to sort
 * @param {number} column The index of the column to sort
 * @param {boolean} asc Determines if the sorting will be in ascending
 */
function sortTableByColumn(table, column, asc = true) {
	const dirModifier = asc ? 1 : -1;
	const tBody = table.tBodies[0];
	const rows = Array.from(tBody.querySelectorAll("tr"));

	// Sort each row
	const sortedRows = rows.sort((a, b) => {
		let aColText = a.querySelector(`td:nth-child(${column + 1})`).textContent.trim();
		let bColText = b.querySelector(`td:nth-child(${column + 1})`).textContent.trim();

		let splittedA = aColText.split(".")

		if (!isNaN(aColText)) {
			aColText = parseFloat(aColText)
			bColText = parseFloat(bColText)

		} else if (!isNaN(new Date(splittedA[2], splittedA[1], splittedA[0]).getDate())) {
			aColText = new Date(splittedA[2], splittedA[1], splittedA[0])
			let splittedB = bColText.split(".")
			bColText = new Date(splittedB[2], splittedB[1], splittedB[0])
		}
		return aColText > bColText ? (1 * dirModifier) : (-1 * dirModifier);
	});

	// Remove all existing TRs from the table
	while (tBody.firstChild) {
		tBody.removeChild(tBody.firstChild);
	}

	// Re-add the newly sorted rows
	tBody.append(...sortedRows);

	// Remember how the column is currently sorted
	table.querySelectorAll("th").forEach(th => th.classList.remove("th-sort-asc", "th-sort-desc"));
	table.querySelector(`th:nth-child(${column + 1})`).classList.toggle("th-sort-asc", asc);
	table.querySelector(`th:nth-child(${column + 1})`).classList.toggle("th-sort-desc", !asc);
}

document.querySelectorAll(".table-sortable th").forEach(headerCell => {
	headerCell.addEventListener("click", () => {
		const tableElement = headerCell.parentElement.parentElement.parentElement;
		const headerIndex = Array.prototype.indexOf.call(headerCell.parentElement.children, headerCell);
		const currentIsAscending = headerCell.classList.contains("th-sort-asc");

		sortTableByColumn(tableElement, headerIndex, !currentIsAscending);
	});
});

// delete all button
document.getElementById("delete-all-btn").onclick = function () {
	if (confirm("Are you sure, that you want to delete the entire table?")) {
		let text = prompt("Please type the deletion key in the text box to confirm the deletion of the entire table: ")

		// send deletion request
		let xhr = new XMLHttpRequest()
		xhr.open("POST", "/deleteall", true)
		xhr.setRequestHeader("Content-Type", "application/json")
		xhr.send(JSON.stringify({Text: text}))

		// remove table on frontend
		document.getElementById("main-table").remove()
		alert("Reload the page to see if the deletion was successful and to create a new table in this case.")

	}
}

