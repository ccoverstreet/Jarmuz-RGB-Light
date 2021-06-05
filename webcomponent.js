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
			<path d="M60,240 A60,60,0,0,1,120,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,180 A120,120,0,0,1,180,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,120 A180,180,0,0,1,240,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,60 A240,240,0,0,1,300,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
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
