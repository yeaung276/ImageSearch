const INPUT_IMAGE_WIDTH = 125;
let cropper;
let model;

class App {
  constructor() {
    this.image = document.getElementById("search-img");
    this.file = document.getElementById("file");
    this.searchBtn = document.getElementById("search");
    this.loader = document.getElementById("loading");

    this.__loadModel();

    this.file.addEventListener("change", (e) => this.onFileSelect(e));
    this.searchBtn.addEventListener("click", () => {
      this.onSearchClick();
    });
  }

  __cropper = () => {
    this.cropper = new Cropper(this.image, {
      aspectRatio: 1,
    });
  };

  async __loadModel() {
    this.showLoading();
    try {
      this.model = await tf.loadGraphModel("../jsmodel/model.json");
      this.setError("");
    } catch (e) {
      this.setError(e);
      console.error(e);
    } finally {
      this.hideLoading();
    }
  }

  showLoading() {
    this.loader.style.display = "flex";
  }

  hideLoading() {
    this.loader.style.display = "none";
  }

  setError(message) {
    document.getElementById("err-msg").innerText = message;
  }

  startCropper(image) {
    this.image.style.maxWidth = "100%";
    this.image.addEventListener("load", this.__cropper);
    this.image.src = image;
  }

  stopCropper() {
    this.cropper.destroy();
    this.cropper = null;
    this.image.removeEventListener("load", this.__cropper);
    this.image.style.maxWidth = "300px";
  }

  async __getEncoding(canvas) {
    const tensor = tf.browser.fromPixels(canvas);
    const pTensor = tensor.expandDims(0).div(125.5).sub(1);
    const result = this.model.predict(pTensor);
    const encoding = await result.array();
    tensor.dispose();
    pTensor.dispose();
    result.dispose();
    return encoding;
  }

  async __searchImage(encoding) {
    return fetch("http://localhost:50052/v1/search", {
      method: "POST",
      body: JSON.stringify({
        embeddings: encoding[0],
      })
    })
    .then((e) => e.json())
    .then((r) => {
        console.log(r)
        document.getElementById("res-count").textContent = r.resultCount
        document.getElementById("search-result").innerHTML = r.result.map((x) => {
            return `<div><img width=300 src="http://localhost${x.url}"/><div>Similarity Score: ${x.similarity}</div></div>`
        }).join('')
    });
  }

  async __uploadImage(encoding) {
    const files = document.getElementById("file").files;
    if (files.length) {
      const form = new FormData();
      form.append("attachment", files[0]);
      const result = await fetch("http://localhost:50052/uploadFile", {
        method: "post",
        body: form,
      }).then((r) => r.json());
      return fetch("http://localhost:50052/v1/add", {
        method: "post",
        body: JSON.stringify({
          imageUrl: result.FileUrl,
          embeddings: encoding[0],
        }),
      })
        .then(() => (document.getElementById("upload-result").html = "pass"))
        .catch(() => (document.getElementById("upload-result").html = "fail"));
    }
  }

  onFileSelect(e) {
    if (e.target.files) {
      this.startCropper(URL.createObjectURL(e.target.files[0]));
      this.file.disabled = true;
    }
  }

  async onSearchClick() {
    this.showLoading();
    const canvas = this.cropper?.getCroppedCanvas({ width: INPUT_IMAGE_WIDTH });
    this.stopCropper();
    const dataURL = canvas.toDataURL();
    document.getElementById("search-img").src = dataURL;
    document.getElementById("file").disabled = false;
    try {
      const encoding = await this.__getEncoding(canvas);
      await (TAB === "Search" ? this.__searchImage : this.__uploadImage)(
        encoding
      );
      this.setError("");
    } catch (e) {
      console.error(e);
      this.setError(e);
    } finally {
      this.hideLoading();
    }
  }
}

const app = new App();
