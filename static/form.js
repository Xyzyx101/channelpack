(function () {
    "use strict";
    $(document).ready(function () {
        //$("#add-input-btn").click(addNewFileInput);
        //connectInputControls();
        $("#file-upload").change(function () {
            uploadImage();
        })
    });

    function uploadImage() {
        let file = $("#file-upload");
        if (file.val()) {
            file.trigger("submit");
        }
    }

    // function addNewFileInput() {
    //     console.log("new");
    // }

    // function connectInputControls() {
    //     let images = $(".input-image");
    //     images.each(function (index) {
    //         let label = $(this).find(".image-label");
    //         let displayImage = label.find("img");
    //         let imageFile = $(this).find(":file");
    //         let deleteBtn = $(this).find(".list-change-btn");

    //         replaceAttrIndex(label, "for", index);
    //         replaceAttrIndex(imageFile, "id", index);
    //         replaceAttrIndex(imageFile, "name", index);

    //         let thisImageInput = $(this);
    //         deleteBtn.click(function () {
    //             deleteInput(thisImageInput);
    //         });

    //         imageFile.change(function () {
    //             updateImage(thisImageInput);
    //         })

    //         //one delete

    //         console.log(index + ": " + $(this).text());
    //     });
    // }

    // function replaceAttrIndex(element, attr, newIndex) {
    //     let oldValue = element.attr(attr);
    //     let splitIdx = oldValue.indexOf('_');
    //     let newValue = oldValue.slice(0, splitIdx + 1) + newIndex.toString();
    //     element.attr(attr, newValue);
    //     return;
    // }

    // function deleteInput(element) {
    //     let images = $(".input-image");
    //     if (images.length == 1) {
    //         return;
    //     }
    //     element.remove();
    // }

    // function updateImage(element) {
    //     let displayImage = element.find("img");
    //     let imageFile = element.find(":file");
    //     displayImage.attr("src", imageFile.val());
    // }
})();