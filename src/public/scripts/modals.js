function htmlDecode(input) {
	var doc = new DOMParser().parseFromString(input, "text/html");
	return doc.documentElement.textContent;
}

function setDeleteEditFuncs() {
	rows = document.getElementsByTagName("tr")

	// open modal and fill in current values when clicking on the edit button
	for (let i = 1; i < rows.length; i++) {
		rows[i].getElementsByTagName("td")[rows[i].getElementsByTagName("td").length - 2].onclick = function () {
			editModal.style.display = "block"
			let s = rows[i].getElementsByTagName("td")[4].innerHTML.split(".")
			let date = new Date(s[2], parseInt(s[1]) - 1, parseInt(s[0]) + 1)

			document.getElementById("editName").value = htmlDecode(rows[i].getElementsByTagName("td")[1].innerHTML)
			document.getElementById("editDescription").value = htmlDecode(rows[i].getElementsByTagName("td")[2].innerHTML)
			document.getElementById("editCount").value = htmlDecode(rows[i].getElementsByTagName("td")[3].innerHTML)
			document.getElementById("editDate").value = date.toISOString().substring(0, 10)
			document.getElementById("index").value = i
		}
	}

	// update the current values
	document.getElementById("edit-form").onsubmit = function () {
		let i = document.getElementById("index").value

		// edit the entered data in the HTML table
		let row = rows[i]
		let splittedDate = document.getElementById("editDate").value.split("-")
		row.getElementsByTagName("td")[1].innerHTML = document.getElementById("editName").value
		row.getElementsByTagName("td")[2].innerHTML = document.getElementById("editDescription").value
		row.getElementsByTagName("td")[3].innerHTML = document.getElementById("editCount").value
		row.getElementsByTagName("td")[4].innerHTML = splittedDate[2] + "." + splittedDate[1] + "." + splittedDate[0]

		// hide modal
		editModal.style.display = "none"
	}

	// delete items
	let deleteButtons = document.getElementsByClassName("deleteCell")

	for (let i = 0; i < deleteButtons.length; i++) {
		deleteButtons[i].onclick = function () {
			if (confirm("Do you really want to delete this item?")) {
				// push removed row to undocache
				undocache.push(rows[i + 1])

				// delete on frontend
				rows[i + 1].remove()

				// fix right undo element
				console.log(rows[i - 1].getElementsByTagName("td")[1].innerHTML)
				console.log(rows[i].getElementsByTagName("td")[1].innerHTML)
				console.log(rows[i + 1].getElementsByTagName("td")[1].innerHTML)

				// send deletion request
				let text = i + 1
				let xhr = new XMLHttpRequest()
				xhr.open("POST", "/delete", true)
				xhr.setRequestHeader("Content-Type", "application/json")
				xhr.send(JSON.stringify({Text: text.toString()}))

				// decrement following numbers
				for (let i = text; i < rows.length; i++) {
					rows[i].getElementsByTagName("td")[0].innerHTML = (parseInt(rows[i].getElementsByTagName("td")[0].innerHTML) - 1).toString()
				}
				setDeleteEditFuncs()
			}
		}
	}

}


// add, remove and edit item modals
let addModal = document.getElementById("add-form-wrapper")
let deleteModal = document.getElementById("delete-form-wrapper")
let editModal = document.getElementById("edit-form-wrapper")

let rows = document.getElementsByTagName("tr")


// for "add item"

// when the user clicks the button, open the modal 
document.getElementById("add-button").onclick = function () {
	addModal.style.display = "block";
}

// when the user clicks on <span> (x), close the modal
document.getElementsByClassName("close")[0].onclick = function () {
	addModal.style.display = "none";
}

document.getElementById("add-form").onsubmit = function () {
	// add the entered data to the HTML table
	let row = document.getElementById("main-table").insertRow(-1)
	row.insertCell(0).innerHTML = document.getElementsByTagName("tr").length - 1
	row.insertCell(1).innerHTML = document.getElementById("name").value
	row.insertCell(2).innerHTML = document.getElementById("description").value
	row.insertCell(3).innerHTML = document.getElementById("count").value
	let splittedDate = document.getElementById("date").value.split("-")
	row.insertCell(4).innerHTML = splittedDate[2] + "." + splittedDate[1] + "." + splittedDate[0]
	row.insertCell(5).innerHTML = "<i class=\"far fa-edit\"></i>"
	row.insertCell(6).innerHTML = "<i class=\"far fa-trash-alt\"></i>"

	row.getElementsByTagName("td")[5].classList.add("editCell")
	row.getElementsByTagName("td")[6].classList.add("deleteCell")

	// hide modal
	addModal.style.display = "none"
}


// edit buttons
// when the user clicks on <span> (x), close the modal
document.getElementsByClassName("close")[1].onclick = function () {
	editModal.style.display = "none"
}


// escape to close modals
document.addEventListener("keydown", function (ev) {
	if (ev.key == "Escape") {
		addModal.style.display = "none"
		editModal.style.display = "none"
	}
})

