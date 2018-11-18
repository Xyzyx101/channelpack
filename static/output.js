(function () {
    "use strict";
    $(document).ready(function () {
        $(".file").change(function (e) {
            populateChannels($(e.target));
        })
        $("#pack-type").change(displayChannelSections);
        $("#filename").change(fixupFilename);
        $("#file-type").change(fixupFilename);
        $("#file-type").change(fileTypeChange);
        let fileSelects = $(".file");
        $.each(fileSelects, function (_, fileSelect) {
            populateChannels($(fileSelect));
        });
        displayChannelSections();
        outputFileProgress();
    });

    function populateChannels(fileSelect) {
        let fileSelectID = fileSelect.attr("id");
        let color = fileSelectID.slice(0, fileSelectID.indexOf("-"));
        let channelSelect = $("#" + color + "-channel");

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

    function displayChannelSections(e) {
        let packType = $("#pack-type").val();
        let allPackTypes = ConfigData.AllPackTypes;
        var activeImageChannels;
        for (var idx = 0; idx < allPackTypes.length; ++idx) {
            if (allPackTypes[idx].Name == packType) {
                let imageChannels = allPackTypes[idx].ImageChannels
                activeImageChannels = imageChannels.split('|');
                break;
            }
        }
        let allImageChannels = ConfigData.AllChannels;
        allImageChannels.forEach(function (channel) {
            let channelSectionNode = $("#" + channel.Name)
            let fileNode = $("#" + channel.Name + "-file")
            let channelNode = $("#" + channel.Name + "-channel");
            var found = activeImageChannels.find(function (activeChannel) {
                return activeChannel == channel.Name;
            })
            channelSectionNode.attr("hidden", !found);
            fileNode.prop('required', !!found);
            channelNode.prop('required', !!found);
        });
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
        let alphaSection = $("#alpha");
        let alphaFile = $("#alpha-file");
        let alphaChannel = $("#alpha-channel");
        switch (fileType) {
            case "jpg":
                alphaSection.hide();
                alphaFile.prop('required', false);
                alphaChannel.prop('required', false);
                break;
            case "png":
            case "tiff":
                alphaSection.show();
                alphaFile.prop('required', true);
                alphaChannel.prop('required', true);
                break;
        }
    }

    function outputFileProgress() {
        let processBtnSection = $("#process-btn-section")
        let progressBarSection = $("#progress-bar-section");
        let progressBar = progressBarSection.find("progress");
        const url = "/output";
        $.ajax({
            url: url,
            cache: false,
            type: 'GET',
            success: function (data, textStatus, jqXHR) {
                if (jqXHR.status == 204 /*no content*/) {
                    processBtnSection.show();
                    progressBarSection.hide();
                } else if (jqXHR.status == 200 /* ok */) {
                    processBtnSection.hide();
                    progressBarSection.show();
                    var tokens = data.split(':');
                    progressBar.attr("value", tokens[1]);
                    progressBar.attr("max", tokens[2]);
                    outputFileProgress();
                } else if (jqXHR.status == 201 /* resource created */) {
                    console.log("Resource created!!!");
                    processBtnSection.show();
                    progressBarSection.hide();
                    var tokens = data.split(':');
                    var a = $("<a>")
                        .attr("href", tokens[1])
                        .attr("download", tokens[2])
                        .appendTo("body");

                    a[0].click();
                    a.remove();
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.log("error:" + textStatus + " " + errorThrown);
            }
        })
    }
})();