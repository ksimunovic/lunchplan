window.addEventListener('load', function () {
    $('#addEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        let formData = objectifyForm($form.serializeArray());
        formData.tags = $form.find(".tagsinput").not(':first').tagsinput('items');
        let tagsInput = $form.find(".tagsinput").not(':first').parent().find('.tt-input').val() + ";";
        tagsInput = tagsInput.split(';');
        for (i = 0; i < tagsInput.length; i++) {
            let newTag = tagsInput[i].trim();
            if (newTag != "") {
                formData.tags.push({"name": tagsInput[i].trim()});
            }
        }
        $.ajax({
            url: "/api/meal",
            type: 'POST',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
            },
            data: JSON.stringify(formData),
            dataType: "json",
            success: function (data) {
                $.ajax({
                    url: "/api/meal/all",
                    type: 'GET',
                    beforeSend: function (xhr) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    },
                    success: function (data) {
                        loadTable("#meals-table", data);
                        $form.trigger('reset');
                        $modal.modal('hide');

                    },
                    error: function () {
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                });
            },
            error: function () {
                $form.trigger('reset');
                $modal.modal('hide');
            },
        });
    })

    $('#editEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        let formData = objectifyForm($form.serializeArray());
        formData.tags = $form.find(".tagsinput").not(':first').tagsinput('items');
        let tagsInput = $form.find(".tagsinput").not(':first').parent().find('.tt-input').val() + ";";
        tagsInput = tagsInput.split(';');
        for (i = 0; i < tagsInput.length; i++) {
            let newTag = tagsInput[i].trim();
            if (newTag != "") {
                formData.tags.push({"name": tagsInput[i].trim()});
            }
        }
        $.ajax({
            url: "/api/meal/" + $form.find('[name="id"]').val(),
            type: 'POST',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
            },
            data: JSON.stringify(formData),
            dataType: "json",
            success: function (data) {
                $.ajax({
                    url: "/api/meal/all",
                    type: 'GET',
                    beforeSend: function (xhr) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    },
                    success: function (data) {
                        loadTable("#meals-table", data);
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                    error: function () {
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                });
            },
            error: function () {
                $form.trigger('reset');
                $modal.modal('hide');
            },
        });
    })

    $('#deleteEmployeeModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');

        let idsArray = $form.serializeArray();
        $.each(idsArray, function (index, object) {
            $.ajax({
                url: "/api/meal/" + object.value,
                type: 'DELETE',
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
                },
                dataType: "json",
                success: function (data) {
                },
                error: function () {
                },
            });
        })
        $.ajax({
            url: "/api/meal/all",
            type: 'GET',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
            },
            success: function (data) {
                loadTable("#meals-table", data);
                $form.trigger('reset');
                $modal.modal('hide');
            },
            error: function () {
                $form.trigger('reset');
                $modal.modal('hide');
            },
        });
    })

    /// Calendar section starts here

    $('#calendar').fullCalendar({
        height: 500,
        events: eventsJson,
        header: {
            left: 'prev month agendaWeek today ',
            center: 'title',
            right: 'meals_page add_event next'
        },
        eventClick: function (calEvent, jsEvent, view) {
            editEvent(calEvent);
        },
        customButtons: {
            add_event: {
                text: 'Add',
                click: function () {
                    $('#addEventModal').modal('show');
                }
            },
            meals_page: {
                text: 'Meals',
                click: function () {
                    window.location = '/';
                }
            }
        },
        dayClick: function (date, jsEvent, view) {
            $addEventModal = $('#addEventModal');
            $addEventModal.find('#datepicker').datepicker('setDate', date.format());
            $addEventModal.modal('show');
        }
    });

    $('#addEventModal [type="submit"], #editEventModal [type="submit"]').on('click', function (e) {
        e.preventDefault();
        let $modal = $(this).closest('.modal');
        let $form = $modal.find('form');
        let $calendar = $('#calendar');
        let link = '';
        if ($(e.target).closest('div .modal').attr('id') == "editEventModal") {
            link = '/' + $form.find('[name="id"]').val();
        }

        let formData = objectifyForm($form.serializeArray());
        $.ajax({
            url: "/api/calendar" + link,
            type: 'POST',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
            },
            data: JSON.stringify(formData),
            dataType: "json",
            success: function () {
                $.ajax({
                    url: "/api/calendar/all",
                    type: 'GET',
                    beforeSend: function (xhr) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                    },
                    success: function (data) {
                        eventsJson = data;
                        $calendar.fullCalendar('removeEventSources');
                        $calendar.fullCalendar('addEventSource', eventsJson);
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                    error: function () {
                        $form.trigger('reset');
                        $modal.modal('hide');
                    },
                });
            },
            error: function () {
                $form.trigger('reset');
                $modal.modal('hide');
            },
        });
    });

    $('.datepicker').datepicker({
        autoclose: true
    });

    $('.meal_search').each(function (index, value) {
        let $searchInput = $(value);
        $searchInput.typeahead({
            source: mealsJson,
            displayText: function (item) {
                $searchInput.parent().find(".meal_id").val('');
                return item.title
            },
            afterSelect: function (item) {
                $searchInput.parent().find(".meal_id").val(item.id);
            }
        });
    });
    $('[data-toggle="tooltip"]').tooltip();

    $('.glyphicon-refresh').parent().on('click', function (e) {
        let $this = $(e.target);
        $this.attr('disabled', 'disabled');

        let $mealAc = $this.parent().parent().find('.meal_search');
        let $mealInput = $this.parent().parent().find('[name="meal_id"]');
        $.ajax({
            url: "/api/meal/suggest",
            type: 'GET',
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
            },
            success: function (data) {
                $this.removeAttr('disabled');
                $mealAc.val(data.title);
                $mealInput.val(data.id);
            },
            error: function () {
                $this.removeAttr('disabled');
                $mealAc.val('');
                $mealInput.val('');
            },
        });
    })
});

function editEntity(target) {
    let id = $(target).closest('tr').find('[data-field="id"]').html();
    let $modal = $('#editEmployeeModal');
    $.ajax({
        url: "/api/meal/" + id,
        type: 'GET',
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
        },
        data: {},
        success: function (data) {
            for (var key in data) {
                if ($modal.find('[name="' + key + '"]').length != 0) {
                    $modal.find('[name="' + key + '"]').val(data[key]);
                }
            }
            let $tagsInput = $modal.find('.tagsinput').not(':first');
            $tagsInput.tagsinput('removeAll');
            for (let tagIndex in data.tags) {
                $tagsInput.tagsinput('add', data.tags[tagIndex]);
            }
            $modal.modal('show');
        },
        error: function () {
        },
    });
}

function deleteEntity(target) {
    let $modal = $('#deleteEmployeeModal');
    let $form = $modal.find('form');
    $form.find('input[name="id"]').remove();
    if (target == "checkbox") {
        $('[name="options[]"]:checked').each(function (index, value) {
            let val = $(this).closest('tr').find('[data-field="id"]').html();
            if (val != "") {
                let input = document.createElement("input");
                input.type = "text";
                input.name = "id";
                input.hidden = "hidden";
                input.setAttribute('value', val);
                $form.find('.modal-body')[0].appendChild(input);
            }
        });
    } else {
        let val = $(target).closest('tr').find('[data-field="id"]').html();
        let input = document.createElement("input");
        input.type = "text";
        input.name = "id";
        input.hidden = "hidden";
        input.setAttribute('value', val);
        $form.find('.modal-body')[0].appendChild(input);
        $modal.modal('show');
    }
}

function getCookieValue(a) {
    var b = document.cookie.match('(^|;)\\s*' + a + '\\s*=\\s*([^;]+)');
    return b ? b.pop() : '';
}

function objectifyForm(formArray) {//serialize data function

    var returnArray = {};
    for (var i = 0; i < formArray.length; i++) {
        returnArray[formArray[i]['name']] = formArray[i]['value'];
    }
    return returnArray;
}

function loadTable(table, data) {
    $("tbody tr:not(:first)").remove();
    let $table = $(table);
    let dataArray = data;
    let $firstRow = $table.find('tbody tr:first');

    $.each(dataArray, function (index, object) {
        let $row = $firstRow.clone();
        $table.find('tbody tr:last').after($row);
        $.each(object, function (key, value) {
            let $td = $row.find('[data-field="' + key + '"]');
            if ($td.length != 0) {
                $td.html(value);
            }
            //alert( key + ": " + value );
        });
        $row.show();
    });

    $("#selectAll").unbind('click');
    $("#selectAll").click(function () {
        if (this.checked) {
            checkbox.each(function () {
                this.checked = true;
            });
        } else {
            checkbox.each(function () {
                this.checked = false;
            });
        }
    });
    let checkboxSpan = $('table tbody .custom-checkbox').on('click', function (e) {
        if (e.target.type == "checkbox") {
            return
        }
        let chk = $(this).find('input[type="checkbox"]');
        if (chk.prop('checked') == true) {
            chk.prop("checked", false);
        } else {
            chk.prop("checked", true);
        }
    });
    let checkbox = $('table tbody input[type="checkbox"]');
    checkbox.removeAttr('checked');
    checkbox.click(function () {
        if (!this.checked) {
            $("#selectAll").prop("checked", false);
        }
    });

    getAllUserTags();
}

function getAllUserTags() {
    $.ajax({
        url: "/api/tag/all",
        type: 'GET',
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
        },
        data: {},
        success: function (data) {
            possibleUserTags = data;

            $('.tags-input-template-clone').remove();
            $('.tags-input-template').each(function (index, value) {
                $template = $(value);
                $templateClone = $template.clone();

                $template.after($templateClone);
                $templateClone.show();
                $templateClone.addClass('tags-input-template-clone');

                tags = new Bloodhound({
                    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('name'),
                    queryTokenizer: Bloodhound.tokenizers.whitespace,
                    local: possibleUserTags
                })
                tags.initialize();
                var elt = $templateClone.find('.tagsinput');
                elt.tagsinput({
                    itemValue: 'id',
                    itemText: 'name',
                    typeaheadjs: {
                        name: 'tags',
                        displayKey: 'name',
                        source: tags.ttAdapter(),
                    }
                });
                elt.on('typeahead:selected', function (event, data) {
                    $('.tagsinput').val(data);
                });
                //elt.tagsinput('add', possibleUserTags[0]);

                $(".twitter-typeahead").css('display', 'inline');
            });


        },
        error: function () {
        },
    });
}


function editEvent(target) {
    let id = target.id;
    let $modal = $('#editEventModal');
    $.ajax({
        url: "/api/calendar/" + id,
        type: 'GET',
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
        },
        data: {},
        success: function (data) {
            for (var key in data) {
                if ($modal.find('[name="' + key + '"]').length != 0) {
                    $modal.find('[name="' + key + '"]').val(data[key]);
                }
            }
            $modal.find('#datepicker').datepicker('setDate', data.date.split('T')[0]);
            $modal.find('.meal_id').val(data.meal.id);
            $modal.find('.meal_search').val(data.meal.title);
            $modal.modal('show');
        },
        error: function () {
        },
    });
}

function deleteEvent(target) {
    let $modal = $(target).closest('div .modal');
    let $form = $modal.find('form');
    let id = $form.find('input[name="id"]').val();
    let $calendar = $('#calendar');

    $.ajax({
        url: "/api/calendar/" + id,
        type: 'DELETE',
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
            xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        },
        dataType: "json",
        success: function (data) {
            $.ajax({
                url: "/api/calendar/all",
                type: 'GET',
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getCookieValue("sid"));
                },
                success: function (data) {
                    eventsJson = data;
                    $calendar.fullCalendar('removeEventSources');
                    $calendar.fullCalendar('addEventSource', eventsJson);
                    $form.trigger('reset');
                    $modal.modal('hide');
                },
            });
        },
        error: function () {
        },
    });
}
