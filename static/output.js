(function () {
    "use strict";
    $(document).ready(function () {
        $(".file").change(function (e) {
            populateChannels(e);
        })
        $("#has-alpha").change(onHasAlphaChange);
    });

    function populateChannels(e) {
        let fileSelect = $(e.target)
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

        let selectedIdx = UploadData.Names.findIndex(function (elem) {
            return elem == fileSelect.val();
        })
        if (selectedIdx == -1) {
            populateOptions(channelSelect, "")
        } else {
            populateOptions(channelSelect, UploadData.Channels[selectedIdx])
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
        let hasAlpha = $(e.target);
        console.log(hasAlpha);
        //console.log(ha)
    }
})();