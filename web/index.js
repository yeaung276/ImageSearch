const INPUT_IMAGE_WIDTH = 125;
let cropper;
let model;

class App{
    constructor(){
        this.image = document.getElementById('search-img')
        this.file = document.getElementById('file')
        this.searchBtn = document.getElementById('search')

        this.__loadModel()

        this.file.addEventListener('change',(e) => this.onFileSelect(e))
        this.searchBtn.addEventListener('click', () => this.onSearchClick())
    }

    __cropper = () => {
        this.cropper = new Cropper(this.image, {
            aspectRatio: 1
        });
    }

    async __loadModel(){
        this.showLoading()
        this.model = await tf.loadGraphModel('../jsmodel/model.json');
        this.hideLoading()
    }

    showLoading(){
        document.getElementById('loading').style.display = 'flex'
    }

    hideLoading(){
        document.getElementById('loading').style.display = 'none'
    }

    startCropper(image){
        this.image.style.maxWidth = '100%'
        this.image.addEventListener('load', this.__cropper)
        this.image.src = image
    }

    stopCropper(){
        this.cropper.destroy()
        this.cropper = null
        this.image.removeEventListener('load', this.__cropper)
        this.image.style.maxWidth = '300px'
    }

    async __getEncoding(canvas){
        const tensor = tf.browser.fromPixels(canvas) 
        const pTensor = tensor.expandDims(0).div(125.5).sub(1);
        const result = this.model.predict(pTensor)
        const encoding = await result.array()
        tensor.dispose()
        pTensor.dispose()
        result.dispose()
        return encoding
    }

    async __searchImage(encoding){
        console.log(encoding)
    }

    onFileSelect(e){
        if(e.target.files){
            this.startCropper(URL.createObjectURL(e.target.files[0]))
            this.file.disabled = true
        }
    }

    async onSearchClick(){
        const canvas = this.cropper?.getCroppedCanvas({width: INPUT_IMAGE_WIDTH})
        this.stopCropper()
        const dataURL = canvas.toDataURL()
        document.getElementById('search-img').src = dataURL
        document.getElementById('file').disabled = false
        this.showLoading()
        const encoding = await this.__getEncoding(canvas)
        const result = await this.__searchImage(encoding)
        this.hideLoading()
    }
}

const app = new App()

