// for "add item"
// get the modal
let addModal = document.getElementById("add-form")

// When the user clicks the button, open the modal 
document.getElementById("add-button").onclick = function () {
	addModal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
document.getElementsByClassName("close")[0].onclick = function () {
	addModal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
	if (event.target == addModal) {
		addModal.style.display = "none";
	}
}

document.getElementById("submitAdd").onclick = function () {
	// add the entered data to the HTML table
	let row = document.getElementById("main-table").insertRow(-1)
	row.insertCell(0).innerHTML = document.getElementById("name").value
	row.insertCell(1).innerHTML = document.getElementById("description").value
	row.insertCell(2).innerHTML = document.getElementById("count").value
	row.insertCell(3).innerHTML = document.getElementById("date").value

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

document.getElementById("submitDelete").onclick = function () {
	// delete the row from the HTML table
	document.getElementsByTagName("tr")[parseInt(document.getElementById("number").value)].remove()

	// hide modal
	deleteModal.style.display = "none"

}
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
		const aColText = a.querySelector(`td:nth-child(${column + 1})`).textContent.trim();
		const bColText = b.querySelector(`td:nth-child(${column + 1})`).textContent.trim();

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
