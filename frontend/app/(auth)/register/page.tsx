import AuthForm from "@/components/common/AuthForm";

const Register = () => {
  return (
    <div className="flex flex-col items-center justify-center w-full max-h-full">
      <AuthForm isLogin={false} />
    </div>
  );
};

export default Register;
