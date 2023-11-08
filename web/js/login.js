(async function () {
	// Fetch all the forms we want to apply custom Bootstrap validation styles to
	var forms = document.querySelectorAll('.needs-validation')

	// Loop over them and prevent submission
	Array.prototype.slice.call(forms)
		.forEach(function (form) {
			form.addEventListener('submit', function (event) {
				// if (!form.checkValidity()) {
				event.preventDefault()
				event.stopPropagation()
				var query = document.getElementById("query_string")
				license = queryLicense(query.value)
				if (license != null) {
					query.classList.remove('is-invalid')
					query.classList.add('is-valid')
					setTimeout(function () {
						window.location.href = `invoice.html?activeCode=${query.value}`
					}, 500)
				}
				else {
					query.classList.remove('is-valid')
					query.classList.add('is-invalid')
				}
			}, false)

		})

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