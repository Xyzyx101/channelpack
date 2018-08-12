(function () {
    "use strict";
    $(document).ready(function () {
        $(".file").change(function (e) {
            populateChannels($(e.target));
        })
        $("#has-alpha").change(onHasAlphaChange);
        $("#filename").change(fixupFilename);
        $("#file-type").change(fixupFilename);
        $("#file-type").change(fileTypeChange);
        let fileSelects = $(".file");
        $.each(fileSelects, function (_, fileSelect) {
            populateChannels($(fileSelect));
        });
        onHasAlphaChange();
    });

    function populateChannels(fileSelect) {
        let channelNodes = $(".channel");
        var channelSelect = null;
        if (fileSelect.attr("id") == "red-file") {
            channelSelect = $("#red-channel");
        } else if (fileSelect.attr("id") == "green-file") {
            channelSelect = $("#green-channel");
        } else if (fileSelect.attr("id") == "blue-file") {
            channelSelect = $("#blue-channel");
        } else if (fileSelect.attr("id") == "alpha-file") {
            channelSelect = $("#alpha-channel");
        }

        let selectedIdx = ConfigData.ImageNames.findIndex(function (elem) {
            return elem == fileSelect.val();
        })
        if (selectedIdx == -1) {
            populateOptions(channelSelect, "")
        } else {
            populateOptions(channelSelect, ConfigData.ImageChannels[selectedIdx])
        }
        channelSelect.val(0);
    }

    function populateOptions(selectElem, options) {
        selectElem.find("option").remove();
        let newOptions = options.split('|');
        $.each(newOptions, function (idx, val) {
            selectElem.append($("<option></option>")
                .attr("value", val).text(val));
        });
    }

    function onHasAlphaChange(e) {
        let hasAlpha = $("#has-alpha");
        let alphaSection = $("#alpha-section");
        let alphaFile = $("#alpha-file");
        let alphaChannel = $("#alpha-channel");
        let alphaActive = hasAlpha.prop('checked')
        alphaActive ? alphaSection.show() : alphaSection.hide();
        alphaFile.prop('required', alphaActive);
        alphaChannel.prop('required', alphaActive);
    }

    function fixupFilename(e) {
        var filename = $("#filename").val();
        let filetype = $("#file-type").val();
        let periodIdx = filename.lastIndexOf(".");
        if (periodIdx > 0) {
            filename = filename.slice(0, periodIdx);
        }
        $("#filename").val(filename + "." + filetype);
    }

    function fileTypeChange(e) {
        let fileType = $("#file-type").val();
        let hasAlpha = $("#has-alpha-section");
        let alphaSection = $("#alpha-section");
        let alphaFile = $("#alpha-file");
        let alphaChannel = $("#alpha-channel");
        switch (fileType) {
            case "jpg":
                alphaSection.hide();
                hasAlpha.hide();
                hasAlpha.prop('checked', false);
                alphaFile.prop('required', false);
                alphaChannel.prop('required', false);
                break;
            case "png":
            case "tga":
                alphaSection.show();
                hasAlpha.show();
                hasAlpha.prop('checked', true);
                alphaFile.prop('required', true);
                alphaChannel.prop('required', true);
                break;
        }
    }
})();