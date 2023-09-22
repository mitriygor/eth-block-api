use yew::prelude::*;

#[derive(Clone, PartialEq, Properties)]
pub struct InputFieldProps {
    // pub input_value: String,
    // pub on_cautious_change: Callback<ChangeData>,
    pub label: String,
    pub field_type: String,
    pub name: String,
    pub input_node_ref: NodeRef,
    pub loading: bool,
}

#[function_component(InputField)]
pub fn input_field(props: &InputFieldProps) -> Html {
    let InputFieldProps {
        // input_value,
        // on_cautious_change,
        label,
        field_type,
        name,
        input_node_ref,
        loading,
    } = props;

    html! {
        <div>
            <label for={name.clone()}  class="form-label">{ label }</label>
            <div class="input-group">
                <input class="form-control"
                    /* onchange={on_cautious_change} */
                    type={field_type.clone()}
                    /* value={input_value.clone()} */
                    name={name.clone()}
                    id={name.clone()}
                    ref={input_node_ref.clone()}
                />
                <button type="submit" class="btn btn-outline-secondary" disabled={*loading}>{">"}</button>
            </div>
        </div>
    }
}