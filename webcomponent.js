class extends HTMLElement {
	constructor() {
		super();
		console.log("HELLO");

		this.attachShadow({mode: "open"});
		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
</style>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display:flex">
		<h1>RGB Light</h1>
		<svg viewBox="0 0 360 360">
			<path d="M150,300 A30,60,0,0,1,210,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M105,300 A60,90,0,0,1,255,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,300 A60,80,0,0,1,300,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
		</svg>
	</div>
	<hr>
	<div class="jmod-body">
		<div id="debug-out">
		</div>
		<p>This will have sliders for controlling the color and a selector for which lights to send the set command to</p>
	</div>
</div>
		`
		this.inner
	}

	init(source, instName, config) {
		this.source = source;
		this.instName = instName;
		this.config = config;

		const elem = this.shadowRoot.getElementById("debug-out");
		for (var x of this.config.lightIPs) {
			elem.textContent += "IP:" + x + "\n";
		}
	}
}
