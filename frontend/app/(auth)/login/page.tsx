import AuthForm from "@/components/common/AuthForm";

const Login = () => {
  return (
    <div className="flex flex-col items-center justify-center w-full max-h-full">
      <AuthForm isLogin={true} />
    </div>
  );
};

export default Login;
