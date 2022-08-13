const puppeteer = require('puppeteer');

(async () => {
	if(process.env.WEBSITE == undefined) {
		console.log("No website to screenshot given.");
		return;
	}

	const browser = await puppeteer.launch();
	const page = await browser.newPage();
	page.setViewport({
		width: 425,
		height: 332,
		deviceScaleFactor:1,
	})
	await page.goto(process.env.WEBSITE);
	// wait up to a minute before screenshotting.
	time = Math.ceil((Math.random()*30000)+15000);
	await page.waitForTimeout(time);
	await page.screenshot({path: 'screenshot.png'});
	await browser.close();
})();