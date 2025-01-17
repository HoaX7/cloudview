/** @type {import('tailwindcss').Config} */
export default {
	content: [ "./src/**/*.{html,js,svelte,ts}" ],
	theme: {
		extend: {
			scale: {
				175: "1.75",
				200: "2",
				500: "5"
			},
			backgroundImage: {
				"gradient":
					"linear-gradient(81.02deg, #4d91ff -23.47%, #b14bf4 45.52%, #fa5560 114.8%)"
			}
		}
	},
	plugins: []
};
