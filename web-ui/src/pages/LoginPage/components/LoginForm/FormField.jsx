import Input from '../../../../shared/ui/Input/Input';

function FormField({ label, id, name, type = 'text', placeholder, value, onChange, autoComplete }) {
  return (
    <Input
      label={label}
      id={id}
      name={name}
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      autoComplete={autoComplete}
    />
  );
}

export default FormField;
