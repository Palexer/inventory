let undocache = new Array()

// add/remove item modals
// for "add item"
// get the modal
let addModal = document.getElementById("add-form")

// when the user clicks the button, open the modal 
document.getElementById("add-button").onclick = function () {
	addModal.style.display = "block";
}

// when the user clicks on <span> (x), close the modal
document.getElementsByClassName("close")[0].onclick = function () {
	addModal.style.display = "none";
}

// when the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
	if (event.target == addModal) {
		addModal.style.display = "none";
	}
}

document.getElementById("submitAdd").onclick = function () {
	// add the entered data to the HTML table
	let row = document.getElementById("main-table").insertRow(-1)
	row.insertCell(0).innerHTML = document.getElementsByTagName("tr").length - 1
	row.insertCell(1).innerHTML = document.getElementById("name").value
	row.insertCell(2).innerHTML = document.getElementById("description").value
	row.insertCell(3).innerHTML = document.getElementById("count").value
	let splittedDate = document.getElementById("date").value.split("-")
	row.insertCell(4).innerHTML = splittedDate[2] + "." + splittedDate[1] + "." + splittedDate[0]

	// hide modal
	addModal.style.display = "none"
}


// for "delete item"
// get the modal
let deleteModal = document.getElementById("delete-form")

// When the user clicks the button, open the modal 
document.getElementById("delete-button").onclick = function () {
	deleteModal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
document.getElementsByClassName("close")[1].onclick = function () {
	deleteModal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
	if (event.target == deleteModal) {
		deleteModal.style.display = "none";
	}
}

// delete the row from the HTML table
document.getElementById("submitDelete").onclick = function () {
	let n = parseInt(document.getElementById("number").value)
	if (n > 0) {
		let tds = document.getElementsByTagName("td")

		for (i = 0; i < tds.length; i++) {
			if (parseInt(tds[i].innerHTML) == n) {
				undocache.push(tds[i].closest("tr"))
				tds[i].closest("tr").remove()
				// hide modal
				deleteModal.style.display = "none"
				return
			}
		}
	}
}

// undo button
document.getElementById("undo-button").onclick = function () {
	if (undocache.length < 1) {
		alert("No elements in undo cache")
		return
	}

	if (confirm("Restore? Are you sure?")) {
		// undo request for backend
		let xhr = new XMLHttpRequest()
		xhr.open("POST", "/undo", true)
		xhr.send()

		// add row back on frontend
		let row = document.getElementById("main-table").insertRow(-1)
		for (i = 0; i < document.getElementsByTagName("th").length; i++) {
			row.insertCell(i).innerHTML = undocache[undocache.length - 1].children[i].innerHTML
		}

		// remove last element from cache
		undocache.pop()
	}
}

// print button
document.getElementById("print-button").onclick = function () {window.print()}

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

		} else if (new Date(splittedA[2], splittedA[1], splittedA[0]).toString() != "Invalid Date") {
			aColText = new Date(aColText)
			bColText = new Date(bColText)
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

// escape to close modals
document.addEventListener("keydown", function (ev) {
	if (ev.key == "Escape") {
		document.getElementById("add-form").style.display = "none"
		document.getElementById("delete-form").style.display = "none"
	}
})

// delete all button
document.getElementById("delete-all-btn").onclick = function () {
	if (confirm("Are you sure, that you want to delete the entire table?")) {
		let text = prompt("Please type 'Inventory' in the text box to confirm the deletion of the entire table: ")
		if (text == "Inventory") {
			// send deletion request
			let xhr = new XMLHttpRequest()
			xhr.open("POST", "/deleteall", true)
			xhr.setRequestHeader("Content-Type", "application/json")
			xhr.send(JSON.stringify({Text: text}))

			// remove table on frontend
			document.getElementById("main-table").remove()
			alert("Deletion successful. Reload the page to create a new table.")

		} else {
			alert("Failed to delete table: confirmation failed")
		}
	}
}

